package tests

import (
	"testing"
	"time"

	"github.com/JayJoshi500/golang-gorm-app/internal/models"
	"github.com/JayJoshi500/golang-gorm-app/internal/services"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// newTestDB spins up a fresh in-memory sqlite database per test so tests
// stay isolated and don't require a running Postgres instance.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))
	return db
}

func newTestService(t *testing.T) services.AuthService {
	t.Helper()
	db := newTestDB(t)
	return services.NewAuthService(db, "test-secret", time.Hour)
}

func TestRegister_Success(t *testing.T) {
	svc := newTestService(t)

	user, err := svc.Register("Ada Lovelace", "ada@example.com", "supersecret1")

	require.NoError(t, err)
	assert.Equal(t, "ada@example.com", user.Email)
	assert.NotEmpty(t, user.ID)
	assert.NotEqual(t, "supersecret1", user.Password, "password must be hashed, not stored raw")
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.Register("Ada Lovelace", "ada@example.com", "supersecret1")
	require.NoError(t, err)

	_, err = svc.Register("Someone Else", "ada@example.com", "anotherpass1")

	assert.ErrorIs(t, err, services.ErrEmailTaken)
}

func TestLogin_Success(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.Register("Ada Lovelace", "ada@example.com", "supersecret1")
	require.NoError(t, err)

	token, user, err := svc.Login("ada@example.com", "supersecret1")

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "ada@example.com", user.Email)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.Register("Ada Lovelace", "ada@example.com", "supersecret1")
	require.NoError(t, err)

	_, _, err = svc.Login("ada@example.com", "wrong-password")

	assert.ErrorIs(t, err, services.ErrInvalidCredential)
}

func TestLogin_UnknownEmail(t *testing.T) {
	svc := newTestService(t)

	_, _, err := svc.Login("nobody@example.com", "whatever1")

	assert.ErrorIs(t, err, services.ErrInvalidCredential)
}

func TestLogout_BlacklistsToken(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.Register("Ada Lovelace", "ada@example.com", "supersecret1")
	require.NoError(t, err)
	token, _, err := svc.Login("ada@example.com", "supersecret1")
	require.NoError(t, err)

	assert.False(t, svc.IsBlacklisted(token))

	err = svc.Logout(token)

	require.NoError(t, err)
	assert.True(t, svc.IsBlacklisted(token))
}

func TestLogout_InvalidToken(t *testing.T) {
	svc := newTestService(t)

	err := svc.Logout("not-a-real-token")

	assert.ErrorIs(t, err, services.ErrInvalidToken)
}
