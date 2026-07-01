# Conelli Admin Backend

Gin and GORM backend service for the Conelli admin app.

The project follows the DevCamper-style layout:

- `api` wires the HTTP server, middleware, routes, handlers, and store.
- `cmd` contains command-line entry points such as migrations.
- `config` loads environment configuration.
- `devops` contains Docker Compose and deployment configuration.
- `internal/handlers` contains HTTP handlers grouped by domain.
- `internal/store` contains GORM database setup, migrations, DAOs, and repositories.
- `scripts` contains operational helper scripts.

## Run Locally

```sh
cp .env.example .env
make db-create
make migrate up
make run
```

Required environment:

```sh
PORT=8000
HOST=localhost
USER=postgres
PASSWORD=postgres
DB_NAME=conelli_admin
DB_PORT=5432
SSLMODE=disable
CORS_ORIGIN=http://localhost:5173
ADMIN_EMAIL=admin@conelliengineering.com
ADMIN_NAME=Conelli Admin
ADMIN_PASSWORD=change-me
```

`DATABASE_URL` is also supported and takes precedence over the split database settings.

`ADMIN_EMAIL` and `ADMIN_PASSWORD` are the environment-backed admin login credentials.

## Migrations

```sh
make migrate up
make migrate up 1
```

The admin console data store uses GORM with Postgres. On startup the API auto-migrates the admin data table and seeds the default staff, properties, investors, updates, milestones, materials, and documents when the table is empty.
