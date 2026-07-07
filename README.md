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
API_PORT=8000
APP_ENV=development
HOST=localhost
DB_USER=postgres
PASSWORD=postgres
DB_NAME=conelli_admin
DB_PORT=55432
SSLMODE=disable
CORS_ORIGIN=http://localhost:5173,http://localhost:5174,http://localhost:5175
ADMIN_EMAIL=admin@conelliengineering.com
ADMIN_NAME=Conelli Admin
ADMIN_PASSWORD=change-me
AWS_REGION=eu-west-1
AWS_S3_BUCKET=
AWS_S3_PREFIX=conelli/dev
AWS_S3_PUBLIC_URL=
```

`DATABASE_URL` is also supported and takes precedence over the split database settings.

`ADMIN_EMAIL` and `ADMIN_PASSWORD` are the environment-backed admin login credentials.

## Admin CRUD API

`/admin/data` remains available for full-dataset sync. The admin UI can also call resource endpoints that persist into the same backend dataset:

- `GET /admin/properties`
- `POST /admin/properties`
- `GET /admin/properties/:id`
- `PUT /admin/properties/:id`
- `DELETE /admin/properties/:id`
- `GET /admin/investors`
- `POST /admin/investors`
- `GET /admin/investors/:id`
- `PUT /admin/investors/:id`
- `DELETE /admin/investors/:id`
- `PUT /admin/investors/:id/properties/:propertyId`
- `DELETE /admin/investors/:id/properties/:propertyId`
- `GET /admin/users`
- `POST /admin/users`
- `GET /admin/users/:id`
- `PUT /admin/users/:id`
- `DELETE /admin/users/:id`
- `POST /admin/updates`
- `PUT /admin/updates/:id`
- `DELETE /admin/updates/:id`
- `POST /admin/milestones`
- `PUT /admin/milestones/:id`
- `DELETE /admin/milestones/:id`
- `POST /admin/documents`
- `PUT /admin/documents/:id`
- `DELETE /admin/documents/:id`

Property records include public project fields used by Corneli and CIS: `client`, `year`, `tags`, `galleryImages`, `publicDescription`, and `publicOverview`.

## S3 Uploads

Images and documents should be uploaded directly to S3 with a presigned URL:

```http
POST /admin/uploads/presign
Content-Type: application/json

{
  "fileName": "cover.jpg",
  "contentType": "image/jpeg",
  "folder": "properties"
}
```

The response includes `method`, `uploadUrl`, `fileUrl`, `key`, and `expiresIn`. Upload the file with `PUT uploadUrl` using the same content type, then save `fileUrl` in the relevant property/document record.

## Docker Database

Start Postgres with Docker:

```sh
make dev-db-up
```

This uses Docker Compose project `conelli` and starts the database container/service as `conelli-postgres`.
By default Postgres is exposed on host port `55432` to avoid conflicts with other local Postgres containers.
The API container is exposed on `API_PORT`, defaulting to `8000`.

## Docker Environments

Development:

```sh
cp .env.development.example .env
make dev-up
make dev-logs
```

Production behind Nginx:

```sh
cp .env.production.example .env
make prod-up
make prod-logs
```

`devops/docker-compose.prod.yml` runs the API behind an Nginx reverse proxy using `devops/nginx.conf`. The proxy forwards `/health`, `/auth`, and `/admin` to the backend and exposes `/healthz` for proxy health checks.

Production startup validates required values and fails fast when any of these are missing:

```sh
APP_ENV=production
PORT=8000
DATABASE_URL=postgres://...
CORS_ORIGIN=https://admin.conelliengineering.com,https://cis.conelliengineering.com
ADMIN_EMAIL=admin@conelliengineering.com
ADMIN_PASSWORD=replace-with-strong-password
AWS_REGION=eu-west-1
AWS_S3_BUCKET=conelli-prod-assets
```

Keep `PORT=8000` in production unless `devops/nginx.conf` is updated too, because the Nginx upstream and Docker health checks target the API container on port `8000`.

The production image runs as a non-root user, includes a container health check at `/health`, and the Go server handles `SIGTERM` with a graceful shutdown window.

## Migrations

```sh
make migrate up
make migrate up 1
```

The admin console data store uses GORM with Postgres. On startup the API auto-migrates the admin data table and seeds the default staff, properties, investors, updates, milestones, materials, and documents when the table is empty.
