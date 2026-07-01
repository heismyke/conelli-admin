package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.up.sql
var migrationFiles embed.FS

type migration struct {
	number int
	name   string
	path   string
}

func Migrate(ctx context.Context, database *sql.DB, targetNumber int) error {
	if database == nil {
		return fmt.Errorf("database is nil")
	}

	if err := ensureMigrationTable(ctx, database); err != nil {
		return err
	}

	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	applied, err := appliedMigrations(ctx, database)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if targetNumber > 0 && migration.number > targetNumber {
			break
		}

		if applied[migration.number] {
			continue
		}

		if err := applyMigration(ctx, database, migration); err != nil {
			return err
		}
	}

	return nil
}

func ensureMigrationTable(ctx context.Context, database *sql.DB) error {
	_, err := database.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
    number integer PRIMARY KEY,
    name text NOT NULL,
    applied_at timestamptz NOT NULL DEFAULT now()
);`)
	if err != nil {
		return fmt.Errorf("ensure schema_migrations table: %w", err)
	}

	return nil
}

func loadMigrations() ([]migration, error) {
	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return nil, fmt.Errorf("read migrations: %w", err)
	}

	migrations := make([]migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}

		number, err := migrationNumber(entry.Name())
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration{
			number: number,
			name:   entry.Name(),
			path:   filepath.Join("migrations", entry.Name()),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].number < migrations[j].number
	})

	return migrations, nil
}

func migrationNumber(name string) (int, error) {
	prefix, _, ok := strings.Cut(name, "_")
	if !ok {
		prefix = strings.TrimSuffix(name, ".up.sql")
	}

	number, err := strconv.Atoi(prefix)
	if err != nil {
		return 0, fmt.Errorf("invalid migration number in %q: %w", name, err)
	}

	return number, nil
}

func appliedMigrations(ctx context.Context, database *sql.DB) (map[int]bool, error) {
	rows, err := database.QueryContext(ctx, "SELECT number FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var number int
		if err := rows.Scan(&number); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}

		applied[number] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}

	return applied, nil
}

func applyMigration(ctx context.Context, database *sql.DB, migration migration) error {
	statement, err := migrationFiles.ReadFile(migration.path)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", migration.name, err)
	}

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", migration.name, err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, string(statement)); err != nil {
		return fmt.Errorf("apply migration %s: %w", migration.name, err)
	}

	if _, err := tx.ExecContext(
		ctx,
		"INSERT INTO schema_migrations (number, name) VALUES ($1, $2)",
		migration.number,
		migration.name,
	); err != nil {
		return fmt.Errorf("record migration %s: %w", migration.name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", migration.name, err)
	}

	return nil
}
