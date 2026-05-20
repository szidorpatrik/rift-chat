# Rift Chat

A next-generation, fully open-source chat application built for the modern web.

Rift Chat is designed from the ground up for high performance, strict type safety, and enterprise-grade security, featuring a robust **Go** backend and a lightning-fast **Vue 3** frontend.

---

## System Architecture & Tech Stack

Rift Chat uses a modern, decoupled monorepo architecture.

**Backend (`/backend`)**

* **Language:** Go 1.22+
* **Routing:** Native `net/http` ServeMux
* **Authentication:** Stateless, HTTP-Only JWTs with secure long-lived refresh token rotation.
* **Database Toolkit:** `sqlc` (Type-safe SQL generation) & `pgx/v5`
* **Migrations:** `goose` (Embedded and run automatically on startup)

**Frontend (`/frontend`)**

* **Framework:** Vue 3 (Composition API) & Vite
* **State Management:** Pinia
* **Styling & Components:** TailwindCSS & Nuxt UI
* **API Client:** Custom generic fetch wrapper with automatic token refresh interception.

**Infrastructure (`docker-compose.yml`)**

* **Database:** PostgreSQL 18
* **`(WIP)` Caching/PubSub:**  PostgreSQL (using LISTEN/NOTIFY)

---

## Running the Application

### 1. Clone the Repository

```bash
git clone https://github.com/szidorpatrik/rift-chat.git && cd rift-chat
```

### 2. Environment Configuration

Create the necessary environment files and directories:

```bash
cp .env-example .env
mkdir -p ./backend/certs
mkdir -p ./postgres/data
```

Look inside `.env` and update the variables as needed. The most important one to change is `BACKEND_JWT_SECRET` and `POSTGRES_PASSWORD` for security reasons. You can generate a random JWT secret with the following command:

```bash
openssl rand -base64 64
```

```env
#...
POSTGRES_PASSWORD=secretpassword # <- CHANGE THIS! No quotes around passwords!
#...
BACKEND_JWT_SECRET=""
#...
```

Copy your TLS certificates into the `./backend/certs` directory.

Update `.env` with the certificate file names like this:

```env
#...
BACKEND_CERT_NAME="cert.pem" # <- CHANGE THIS to your cert file name
BACKEND_KEY_NAME="key.pem"   # <- CHANGE THIS to your key file name
#...
```

### 4. Start the Application

Run this command to build the Docker images, start the containers, and show the backend logs:

```bash
docker compose build && docker compose up -d --force-recreate && docker compose logs -f backend
```

If any environment variables are missing or misconfigured, the backend will log an error and exit.
Fix the issues and restart the containers.

### 5. Access the Frontend

Open your browser and navigate to `https://localhost/` to see the frontend.

### Production Deployment

When deploying to production, remove all the ports from [compose.yml](docker-compose.yml) except for caddy and set the appropriate environment variables in `.env` to allow the backend to listen on

## Local Development

See [`backend`](backend/README.md) and [`frontend`](frontend/README.md) for detailed instructions.

## Documentation

* **[API Documentation](docs/api.md):** Detailed specs on all REST endpoints, request/response bodies, and authentication flows.
* **[Backend Documentation](backend/README.md):** In-depth details on the Go architecture, context handling, and database generation.
* **`(WIP)` [Frontend Documentation](frontend/README.md):** Information on the Vue component structure, Pinia state, and routing.
