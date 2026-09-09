package services

import (
	"errors"
	"sync"
	"time"

	"github.com/JayJoshi500/golang-gorm-app/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailTaken        = errors.New("email is already registered")
	ErrInvalidCredential = errors.New("invalid email or password")
	ErrInvalidToken      = errors.New("invalid or expired token")
)

// Claims is the JWT payload issued on login.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// AuthService defines the auth use cases exposed to handlers. Coding to
// an interface (rather than the concrete struct) keeps handlers testable
// with a fake/mock implementation.
type AuthService interface {
	Register(name, email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
	Logout(token string) error
	ParseToken(tokenString string) (*Claims, error)
	IsBlacklisted(tokenString string) bool
}

type authService struct {
	db         *gorm.DB
	jwtSecret  []byte
	jwtExpiry  time.Duration
	blacklist  map[string]time.Time // token -> expiry, for logout on stateless JWTs
	blacklistM sync.Mutex
}

// NewAuthService wires a GORM handle plus JWT settings into an AuthService.
func NewAuthService(db *gorm.DB, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &authService{
		db:        db,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: jwtExpiry,
		blacklist: make(map[string]time.Time),
	}
}

func (s *authService) Register(name, email, password string) (*models.User, error) {
	var existing models.User
	err := s.db.Where("email = ?", email).First(&existing).Error
	if err == nil {
		return nil, ErrEmailTaken
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredential
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredential
	}

	token, err := s.generateToken(&user)
	if err != nil {
		return "", nil, err
	}
	return token, &user, nil
}

// Logout blacklists the token for the remainder of its natural lifetime.
// This is an in-memory demo implementation — swap for Redis (or similar
// shared store) before running more than one instance of the app.
func (s *authService) Logout(tokenString string) error {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return err
	}

	s.blacklistM.Lock()
	defer s.blacklistM.Unlock()
	s.blacklist[tokenString] = claims.ExpiresAt.Time

	s.pruneExpiredLocked()
	return nil
}

// IsBlacklisted reports whether a token was explicitly logged out.
func (s *authService) IsBlacklisted(tokenString string) bool {
	s.blacklistM.Lock()
	defer s.blacklistM.Unlock()
	_, found := s.blacklist[tokenString]
	return found
}

// pruneExpiredLocked drops blacklist entries whose token has already
// expired naturally, since they no longer need tracking. Caller must
// hold blacklistM.
func (s *authService) pruneExpiredLocked() {
	now := time.Now()
	for token, exp := range s.blacklist {
		if now.After(exp) {
			delete(s.blacklist, token)
		}
	}
}

func (s *authService) generateToken(user *models.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpiry)),
			Subject:   user.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ParseToken validates signature and expiry and returns the claims.
func (s *authService) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
