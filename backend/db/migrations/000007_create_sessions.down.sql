-- 000007_create_sessions.down.sql
DROP TRIGGER IF EXISTS trg_sessions_set_updated_at ON sessions;
DROP TABLE IF EXISTS sessions;
