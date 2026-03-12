-- 000002_create_devices.up.sql
CREATE TABLE IF NOT EXISTS devices (
    device_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    push_token   TEXT,
    device_model VARCHAR(100),
    last_active  TIMESTAMPTZ,
    is_revoked   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_devices_user_id ON devices(user_id);
