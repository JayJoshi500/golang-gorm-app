# golang-gorm-app

A minimal, production-shaped REST API skeleton in Go: [go-chi](https://github.com/go-chi/chi) for
routing, [GORM](https://gorm.io) + PostgreSQL for persistence, JWT-based auth, structured logging,
and a unit test setup — ready to build on.

## Features

- **Routing:** go-chi v5, versioned under `/api/v1`
- **Database:** PostgreSQL via GORM, connection pooling, auto-migration
- **Auth:** register / login / logout with bcrypt password hashing + JWT access tokens
  (logout uses an in-memory token blacklist — see note below)
- **Middleware:** panic recovery with clean JSON errors, structured request logging, CORS, request ID, timeout
- **Logging:** structured JSON logs via the standard library's `log/slog`
- **Testing:** `testify` + an in-memory SQLite DB, so unit tests don't need a real Postgres instance
- **Makefile:** common dev commands (`run`, `build`, `test`, `fmt`, `vet`, `tidy`)

## Project structure

```
.
├── cmd/api/main.go              # entrypoint: wiring + graceful shutdown
├── internal/
│   ├── config/                  # env-based configuration
│   ├── database/                # GORM/Postgres connection + migrations
│   ├── models/                  # User entity
│   ├── services/                # auth business logic (register/login/logout, JWT)
│   ├── handlers/                # HTTP handlers (request/response glue)
│   ├── middleware/               # recovery, logging, JWT auth guard
│   └── router/                  # chi route tree
├── pkg/
│   ├── logger/                  # slog setup
│   └── response/                # standard JSON response envelope
├── tests/                       # unit tests (in-memory SQLite)
├── .env.example
├── Makefile
└── go.mod
```

## Prerequisites

- Go 1.22+
- Docker Desktop, for the containerized API and PostgreSQL setup

## Setup

```bash
git clone <this-repo>
cd golang-gorm-app

cp .env.example .env
# edit .env with your DB credentials and a real JWT_SECRET

go mod tidy   # downloads dependencies and generates go.sum
```

> This project was generated offline, so `go.sum` is not included yet — `go mod tidy` will fetch
> the exact dependency versions declared in `go.mod` and lock them.

## Running

```bash
make run
# or: go run ./cmd/api
```

To run the API and PostgreSQL together with Docker:

```bash
docker compose up --build
```

The API is available at `http://localhost:8080`. Set `API_PORT` to expose it on a
different host port. PostgreSQL data is stored in the `postgres-data` Docker volume.
Stop the services with `docker compose down`; add `-v` if you also want to remove the
database volume.

The API starts on `http://localhost:8080` (configurable via `SERVER_PORT`). Tables are
auto-migrated on startup.

## API

All responses use a common envelope: `{"success": bool, "data": ..., "error": "..."}`.

| Method | Path                    | Body                                       | Notes                     |
| ------ | ----------------------- | ------------------------------------------ | ------------------------- |
| GET    | `/health`               | –                                          | Liveness check            |
| POST   | `/api/v1/auth/register` | `{ "name", "email", "password" }`          | Password min 8 chars      |
| POST   | `/api/v1/auth/login`    | `{ "email", "password" }`                  | Returns `{ token, user }` |
| POST   | `/api/v1/auth/logout`   | – (header `Authorization: Bearer <token>`) | Blacklists the token      |

Example:

```bash
curl -X POST localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada Lovelace","email":"ada@example.com","password":"supersecret1"}'

curl -X POST localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ada@example.com","password":"supersecret1"}'

curl -X POST localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer <token-from-login>"
```

A `middleware.RequireAuth` guard is included in `internal/middleware/auth.go`, ready to mount
on a protected route group once you add authenticated endpoints — it isn't applied to the
auth routes themselves since those must stay public.

## Testing

```bash
make test              # go test ./... -v
make test-coverage     # with a coverage summary
```

Tests exercise the auth service against an in-memory SQLite database, so no external services
are required.

## Notes & next steps

- **Logout / token revocation** is implemented with an in-memory map for simplicity. It works for
  a single instance but won't be shared across replicas or survive a restart — swap in Redis (or
  another shared store) before running this in a multi-instance deployment.
- **Migrations** use `AutoMigrate`, which is convenient for development but not a substitute for
  versioned migrations in production (consider `golang-migrate`, `goose`, or `atlas`).
- **Validation** on request bodies is intentionally minimal — consider a library like
  `go-playground/validator` as the API grows.
