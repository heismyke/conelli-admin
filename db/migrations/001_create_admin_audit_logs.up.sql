CREATE TABLE IF NOT EXISTS admin_audit_logs (
	id BIGSERIAL PRIMARY KEY,
	actor_id TEXT,
	action TEXT NOT NULL,
	target_type TEXT,
	target_id TEXT,
	metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
