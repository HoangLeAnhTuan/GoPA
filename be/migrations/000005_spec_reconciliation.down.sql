-- ============================================================================
-- GoPA Migration 000005 (Down): Revert Spec Reconciliation & Restore State
-- ============================================================================

-- Recreate legacy tables if reverting
CREATE TABLE network_nodes (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_name TEXT NOT NULL,
    mac_address TEXT,
    static_ip TEXT,
    vlan_tag SMALLINT,
    connection_type TEXT NOT NULL,
    parent_node_id UUID REFERENCES network_nodes(id) ON DELETE SET NULL,
    location TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE vehicles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    make TEXT NOT NULL,
    model TEXT NOT NULL,
    year SMALLINT NOT NULL,
    nickname TEXT,
    current_mileage INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE vehicle_logs (
    id UUID PRIMARY KEY,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_type TEXT NOT NULL,
    mileage INTEGER NOT NULL,
    description TEXT NOT NULL,
    cost NUMERIC(14, 2) NOT NULL,
    log_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- Drop indexes
DROP INDEX IF EXISTS idx_outbox_pending;
DROP INDEX IF EXISTS idx_vocabularies_review_queue;
DROP INDEX IF EXISTS idx_learning_sessions_user;

-- Revert journals columns
ALTER TABLE journals
    DROP COLUMN IF EXISTS energy_level,
    DROP COLUMN IF EXISTS pinned,
    DROP COLUMN IF EXISTS word_count;

-- Drop learning_sessions table
DROP TABLE IF EXISTS learning_sessions CASCADE;

-- Revert vocabulary_reviews columns
ALTER TABLE vocabulary_reviews
    DROP COLUMN IF EXISTS study_mode,
    DROP COLUMN IF EXISTS response_time_ms,
    DROP COLUMN IF EXISTS was_correct;

-- Revert vocabularies columns
ALTER TABLE vocabularies
    DROP COLUMN IF EXISTS example_translation,
    DROP COLUMN IF EXISTS tags,
    DROP COLUMN IF EXISTS difficulty_level,
    DROP COLUMN IF EXISTS audio_url,
    DROP COLUMN IF EXISTS image_url,
    DROP COLUMN IF EXISTS ease_factor,
    DROP COLUMN IF EXISTS interval_days,
    DROP COLUMN IF EXISTS repetition_count,
    DROP COLUMN IF EXISTS mastery_score;

-- Revert tasks columns and constraints
ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS tasks_priority_check,
    DROP CONSTRAINT IF EXISTS tasks_category_check,
    DROP COLUMN IF EXISTS estimated_pomodoros,
    DROP COLUMN IF EXISTS completed_pomodoros;

ALTER TABLE tasks
    ADD CONSTRAINT tasks_priority_check CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH')),
    ADD CONSTRAINT tasks_category_check CHECK (category IN ('WORK', 'STUDY', 'LIFE'));

-- Revert users columns
ALTER TABLE users
    DROP COLUMN IF EXISTS display_name,
    DROP COLUMN IF EXISTS avatar_url;
