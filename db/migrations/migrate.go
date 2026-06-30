package migrations

import (
	"database/sql"
	"fmt"
)

func Run(db *sql.DB, direction string) error {
	switch direction {
	case "up":
		return runUp(db)
	case "down":
		return runDown(db)
	default:
		return fmt.Errorf("unsupported migration direction: %s", direction)
	}
}

func runUp(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS admin_audit_logs (
			id BIGSERIAL PRIMARY KEY,
			actor_id TEXT,
			action TEXT NOT NULL,
			target_type TEXT,
			target_id TEXT,
			metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

func runDown(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS admin_audit_logs;`)
	return err
}
