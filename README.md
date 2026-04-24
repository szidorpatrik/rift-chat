# RiftChat

A next-generation, fully open-source chat application built for the modern web.

## Development

### Prerequisites

- golang
- air (optional, for hot reloading during development)
- bun
- docker
- docker compose
- sqlc
- goose (optional, for manual migrations only)

### Backend

#### API

To locally run the api:

```sh
cd backend
go run main.go
```

For hot reloading during development, you can use `air`:

Install `air` if you haven't already:

```sh
go install github.com/air-verse/air@latest
```

Then run `air` in the `backend` directory:

```sh
cd backend
air
```

#### Database

To generate the database code using `sqlc`, run the following command in the `backend` directory:

```sh
cd backend
sqlc generate
```

Migrations are managed from code. You can use `goose` to run them manually if needed.

By default it always runs the latest migrations on startup.

### Frontend

To locally run the frontend
Install dependencies and start the development server:

```sh
cd frontend
bun install
bun run dev
```
