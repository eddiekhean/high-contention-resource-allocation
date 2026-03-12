-- 000007_create_sessions.up.sql
CREATE TABLE IF NOT EXISTS sessions (
    session_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    device_id   UUID REFERENCES devices(device_id) ON DELETE SET NULL,
    jti         VARCHAR(255) UNIQUE NOT NULL,
    ip_address  VARCHAR(45),
    user_agent  TEXT,
    expires_at  TIMESTAMPTZ NOT NULL,
    is_revoked  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_jti ON sessions(jti);

CREATE TRIGGER trg_sessions_set_updated_at
    BEFORE UPDATE ON sessions
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
