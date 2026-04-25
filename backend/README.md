# Rift Chat - Backend API

The backend service for **Rift Chat**, a modern, real-time chat application. Built with Go (Golang), this API is designed for high performance, strict type safety, and enterprise-grade security.

## Features

* **Advanced Authentication & Security:**
  * **HTTP-Only Cookies:** JWTs are stored in secure, `SameSite=Strict` HTTP-only cookies to completely prevent XSS attacks.
  * **Token Rotation:** Seamless access token rotation using long-lived refresh tokens.
  * **Atomic Transactions:** Database operations (like token rotation and revocation) are wrapped in strict SQL transactions (`sqlc.WithTx`) to prevent orphaned tokens during partial network failures.
  * **Timing Attack Prevention:** Constant-time password hashing and dummy-hashing fallbacks on invalid usernames.
* **Modern Go Architecture:**
  * Uses Go 1.22+ native multiplexer (`http.NewServeMux`) for semantic route definitions (`POST /api/v1/login`).
  * **Generics:** Implements generic JSON decoding and validation pipelines (`decodeAndValidate[T any]`) to eliminate boilerplate.
  * **Middleware Chain:** Clean separation of concerns using `LoggingMiddleware`, `IdentityMiddleware`, and `RequireAuth`.
* **Structured Logging:** Centralized `slog` implementation with request ID tracking, event categorization, and custom JSON error handling.
* **High Performance Infrastructure:** Powered by PostgreSQL.

## Tech Stack

* **Language:** Go (1.22+)
* **Database:** PostgreSQL
* **Caching / In-Memory:** PostgreSQL
* **Database Toolkit:** `sqlc` (Type-safe SQL generation), `goose` (migrations), `pgx/v5`
* **Live Reloading:** `Air`
* **Infrastructure:** Docker & Docker Compose

## Local Development Workflow

This project uses a **Hybrid Local Setup**. The infrastructure (Postgres) runs in Docker, while the Go code runs natively on your machine using `Air` for instant hot-reloading.

### 1. Prerequisites

* [Go](https://go.dev/) (1.22 or higher)
* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [gopls](https://pkg.go.dev/golang.org/x/tools/gopls) (Go language server for IDE features)
  * `go install golang.org/x/tools/gopls`
* [air](https://github.com/cosmtrek/air) (for live-reloading Go development)
  * `go install github.com/air-verse/air@latest`
* [sqlc](https://sqlc.dev/) (for generating DB queries)
  * `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
* [goose](https://github.com/pressly/goose) (for database migrations)
  * `go install github.com/pressly/goose/v3/cmd/goose@latest`

That makes perfect sense! Your `LoadEnvs` function handles this beautifully by explicitly looking for `.env-dev` when running natively, and relying on Docker's environment injection when in `prod`. This is a highly professional setup.

### 2. Environment Configuration

This project uses a dual-environment file setup to handle the networking differences between running Go natively vs inside Docker.

You need to create two files in the root directory: `.env` (for Docker/Prod) and `.env-dev` (for local Go development).

Create them by copying the provided examples:

```bash
cp .env-example .env
cp backend/.env-dev-example backend/.env-dev
```

**File 1: `.env` (Used by Docker Compose)**
This file is read by Docker. It configures the containers to talk to each other using internal Docker network hostnames.

```env
# --- Container names ---
POSTGRES_CONTAINER_NAME="riftchat_postgres"
BACKEND_CONTAINER_NAME="riftchat_backend"
FRONTEND_CONTAINER_NAME="riftchat_frontend"
CADDY_CONTAINER_NAME="riftchat_caddy"

# --- Database ---
POSTGRES_USER="riftchat_user"
POSTGRES_PASSWORD=production-password # <- CHANGE THIS! No quotes around passwords!
POSTGRES_DB="riftchat_db"
POSTGRES_PORT="5432"
POSTGRES_SSLMODE="disable"

# --- Backend ---
BACKEND_HTTP_PORT="8081"
BACKEND_QUIC_PORT="4242"
BACKEND_PLATFORM="prod"
BACKEND_LISTEN_ADDR="0.0.0.0"
# Generate a random key with `openssl rand -base64 64` and paste it below.
BACKEND_JWT_SECRET=""
BACKEND_DATABASE_URL="host=${POSTGRES_CONTAINER_NAME} port=${POSTGRES_PORT} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=${POSTGRES_SSLMODE}"

# --- Frontend ---
FRONTEND_PORT="80"
```

**File 2: `backend/.env-dev` (Used by Air / Native Go)**
When you run the backend natively using `air`, Go cannot see the Docker internal network. This file overrides the platform and database host to point to `127.0.0.1` or `localhost` so Go can connect to the database exposed on your machine's port.

### 3. Start the Infrastructure

Spin up PostgreSQL in the background. You only need to run this once.

```bash
docker compose up -d db
```

Create a new database for development, change the credentials in `backend/.env-dev` if you want to customize them:

```bash
docker compose exec db psql -U riftchat_user -d riftchat_db -c "CREATE USER riftchat_user_dev WITH PASSWORD 'development-password';" -c "CREATE DATABASE riftchat_db_dev OWNER riftchat_user_dev;"
```

### 4. Run the API

Start the Go server using Air. It will automatically recompile and restart whenever you save a Go file.

```bash
air
```

The API will now be listening on `http://0.0.0.0:8081`.

## Project Structure

```text
rift-chat-api/
├── internal/
│   ├── api/            # HTTP handlers, middleware, and routing logic
│   ├── auth/           # JWT generation, validation, and password hashing
│   ├── database/       # Generated sqlc code (Do not edit manually)
│   ├── logger/         # Structured slog configuration
│   └── models/         # Core domain models and types
├── sql/
│   ├── queries/        # Raw .sql files containing queries
│   └── schema/         # Database migration files (goose)
├── .dockerignore       # Docker ignore file
├── Dockerfile          # Dockerfile for containerizing the API
├── main.go             # Application entrypoint
└── sqlc.yml            # sqlc configuration file
```

## API Documentation

All API endpoints are prefixed with `/api/v1`.

For complete documentation regarding request/response shapes, headers, and status codes, please refer to the [API Documentation](../docs/api.md) file.

**Core Routes:**

* `GET  /healthz` - System health check
* `POST /api/v1/users` - Register a new account
* `POST /api/v1/login` - Authenticate and receive HTTP-Only session cookies
* `GET  /api/v1/session` - Restore session state (Requires Auth)
* `POST /api/v1/refresh` - Rotate expired access tokens
* `POST /api/v1/revoke` - Revoke tokens and log out

## Database Changes (sqlc)

When you need to modify the database schema or add new queries:

1. Write your raw SQL in the `sql/queries` or `sql/schema` directories.
2. Run `sqlc generate` to automatically create the type-safe Go bindings in `internal/database`.
3. Use the newly generated methods attached to `cfg.DbQueries`.
