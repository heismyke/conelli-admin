# Setup

1. Install Go 1.21 or newer.
2. Create a PostgreSQL database for the admin service.
3. Copy `.env.example` to `.env` and set `DATABASE_URL`.
4. Start the API with `go run .`.
5. Optional: run SQL migrations with `go run ./cmd/migrate -direction up`; the GORM-backed admin data table is auto-migrated on startup.
