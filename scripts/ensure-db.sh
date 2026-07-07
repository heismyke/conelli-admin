#!/usr/bin/env sh
set -eu

if [ -f ".env" ]; then
	set -a
	. ./.env
	set +a
fi

DB_HOST="${HOST:-localhost}"
DB_PORT="${DB_PORT:-55432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-conelli_admin}"

if [ -n "${PASSWORD:-}" ]; then
	export PGPASSWORD="$PASSWORD"
fi

exists="$(
	psql \
		-h "$DB_HOST" \
		-p "$DB_PORT" \
		-U "$DB_USER" \
		-d postgres \
		-tAc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'"
)"

if [ "$exists" = "1" ]; then
	echo "database '$DB_NAME' already exists"
	exit 0
fi

createdb \
	-h "$DB_HOST" \
	-p "$DB_PORT" \
	-U "$DB_USER" \
	"$DB_NAME"

echo "created database '$DB_NAME'"
