-- 000006_add_triggers.up.sql
-- ─────────────────────────────────────────────────────────────────────────────
-- 1. Generic function: auto-set updated_at on any UPDATE
-- ─────────────────────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to users (has updated_at column)
CREATE TRIGGER trg_users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();


-- ─────────────────────────────────────────────────────────────────────────────
-- 2. Add updated_at to connections + trigger (track when status changes)
-- ─────────────────────────────────────────────────────────────────────────────
ALTER TABLE connections
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

CREATE TRIGGER trg_connections_set_updated_at
    BEFORE UPDATE ON connections
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();


-- ─────────────────────────────────────────────────────────────────────────────
-- 3. Prevent self-connection (user cannot send friend request to themselves)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION prevent_self_connection()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.requester_id = NEW.addressee_id THEN
        RAISE EXCEPTION 'User cannot connect to themselves (requester_id = addressee_id = %)', NEW.requester_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_connections_no_self
    BEFORE INSERT OR UPDATE ON connections
    FOR EACH ROW
    EXECUTE FUNCTION prevent_self_connection();


-- ─────────────────────────────────────────────────────────────────────────────
-- 4. Add role + left_at to conversation_members (group admin / soft-delete)
-- ─────────────────────────────────────────────────────────────────────────────
ALTER TABLE conversation_members
    ADD COLUMN IF NOT EXISTS role    VARCHAR(10) NOT NULL DEFAULT 'member'
        CHECK (role IN ('member', 'admin')),
    ADD COLUMN IF NOT EXISTS left_at TIMESTAMPTZ;   -- NULL = still in group


-- ─────────────────────────────────────────────────────────────────────────────
-- 5. Prevent duplicate 'pair' conversation between the same two users
--    A pair conversation must have exactly 2 members; enforce before insert.
-- ─────────────────────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION prevent_duplicate_pair_conversation()
RETURNS TRIGGER AS $$
DECLARE
    v_type TEXT;
    v_other_user UUID;
    v_existing UUID;
BEGIN
    -- Only enforce for 'pair' type conversations
    SELECT type INTO v_type FROM conversations WHERE conversation_id = NEW.conversation_id;
    IF v_type <> 'pair' THEN
        RETURN NEW;
    END IF;

    -- Find the other member already in this conversation (if any)
    SELECT user_id INTO v_other_user
    FROM conversation_members
    WHERE conversation_id = NEW.conversation_id
      AND user_id <> NEW.user_id
    LIMIT 1;

    IF v_other_user IS NULL THEN
        -- First member being added — always OK
        RETURN NEW;
    END IF;

    -- Check if a pair conversation already exists between these two users
    SELECT cm1.conversation_id INTO v_existing
    FROM conversation_members cm1
    JOIN conversation_members cm2 ON cm1.conversation_id = cm2.conversation_id
    JOIN conversations c ON c.conversation_id = cm1.conversation_id
    WHERE c.type = 'pair'
      AND cm1.user_id = NEW.user_id
      AND cm2.user_id = v_other_user
      AND cm1.conversation_id <> NEW.conversation_id
    LIMIT 1;

    IF v_existing IS NOT NULL THEN
        RAISE EXCEPTION 'A pair conversation already exists between users % and % (conversation_id: %)',
            NEW.user_id, v_other_user, v_existing;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_no_duplicate_pair_conv
    BEFORE INSERT ON conversation_members
    FOR EACH ROW
    EXECUTE FUNCTION prevent_duplicate_pair_conversation();
