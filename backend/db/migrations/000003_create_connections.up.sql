-- 000003_create_connections.up.sql
CREATE TABLE IF NOT EXISTS connections (
    connection_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requester_id  UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    addressee_id  UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status        VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'accepted', 'rejected', 'blocked')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (requester_id, addressee_id)
);

CREATE INDEX IF NOT EXISTS idx_connections_requester ON connections(requester_id, status);
CREATE INDEX IF NOT EXISTS idx_connections_addressee ON connections(addressee_id, status);
