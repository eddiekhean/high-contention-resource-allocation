-- 000006_add_triggers.down.sql
DROP TRIGGER IF EXISTS trg_no_duplicate_pair_conv ON conversation_members;
DROP FUNCTION IF EXISTS prevent_duplicate_pair_conversation();

ALTER TABLE conversation_members
    DROP COLUMN IF EXISTS left_at,
    DROP COLUMN IF EXISTS role;

DROP TRIGGER IF EXISTS trg_connections_no_self ON connections;
DROP FUNCTION IF EXISTS prevent_self_connection();

DROP TRIGGER IF EXISTS trg_connections_set_updated_at ON connections;
ALTER TABLE connections DROP COLUMN IF EXISTS updated_at;

DROP TRIGGER IF EXISTS trg_users_set_updated_at ON users;
DROP FUNCTION IF EXISTS set_updated_at();
