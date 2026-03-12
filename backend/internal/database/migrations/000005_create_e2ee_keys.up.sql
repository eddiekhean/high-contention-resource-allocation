-- 000005_create_e2ee_keys.up.sql
CREATE TABLE IF NOT EXISTS e2ee_keys (
    key_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    device_id       UUID NOT NULL REFERENCES devices(device_id) ON DELETE CASCADE,
    public_key_data TEXT NOT NULL,
    key_type        VARCHAR(20) NOT NULL CHECK (key_type IN ('identity', 'prekey', 'onetime_prekey')),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_e2ee_keys_device ON e2ee_keys(device_id, key_type, is_active);
CREATE INDEX IF NOT EXISTS idx_e2ee_keys_user   ON e2ee_keys(user_id, is_active);
