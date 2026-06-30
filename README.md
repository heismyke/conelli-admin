# Conelli Admin Backend

Gin-based backend service for the Conelli admin app. The folder structure follows the existing Cofikra backend layout:

- `api` for HTTP handlers
- `config` for environment configuration
- `db` for database clients, DAOs, repositories, and migrations
- `http/rq` and `http/rs` for request/response DTOs
- `internal/api` for route definitions and request wrapping
- `internal/router` for Gin setup
- `logic` for business services

## Run Locally

```sh
cp .env.example .env
go mod tidy
go run .
```

Required environment:

```sh
DATABASE_URL=postgres://postgres:postgres@localhost:5432/conelli_admin?sslmode=disable
HTTP_PORT=8000
IS_DEV=true
CORS_ORIGIN=http://localhost:5173
```

The admin console data store uses GORM with Postgres. On startup the API auto-migrates the admin data table and seeds the default staff, properties, investors, updates, milestones, materials, and documents when the table is empty.

## Migrations

```sh
go run ./cmd/migrate -direction up
go run ./cmd/migrate -direction down
```
