-- ============================================================================
-- GoPA Migration 000005: Master Specification Reconciliation & Legacy Cleanup
-- ============================================================================

-- 1. IAM: Add display_name and avatar_url to users table
ALTER TABLE users
    ADD COLUMN display_name TEXT NOT NULL DEFAULT '',
    ADD COLUMN avatar_url TEXT;

-- 2. Tasks: Reconcile priority & category constraints, add pomodoro counters
ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS tasks_priority_check,
    DROP CONSTRAINT IF EXISTS tasks_category_check;

ALTER TABLE tasks
    ADD CONSTRAINT tasks_priority_check CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT')),
    ADD CONSTRAINT tasks_category_check CHECK (category IN ('WORK', 'STUDY', 'LIFE', 'LEETCODE')),
    ADD COLUMN estimated_pomodoros INTEGER NOT NULL DEFAULT 0 CHECK (estimated_pomodoros >= 0),
    ADD COLUMN completed_pomodoros INTEGER NOT NULL DEFAULT 0 CHECK (completed_pomodoros >= 0);

-- 3. Linguistics: Reconcile vocabularies table with full SM-2 algorithm & metadata
ALTER TABLE vocabularies
    ADD COLUMN example_translation TEXT,
    ADD COLUMN tags TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN difficulty_level TEXT NOT NULL DEFAULT 'BEGINNER' CHECK (difficulty_level IN ('BEGINNER', 'INTERMEDIATE', 'ADVANCED')),
    ADD COLUMN audio_url TEXT,
    ADD COLUMN image_url TEXT,
    ADD COLUMN ease_factor NUMERIC(4, 2) NOT NULL DEFAULT 2.50 CHECK (ease_factor >= 1.30),
    ADD COLUMN interval_days INTEGER NOT NULL DEFAULT 0 CHECK (interval_days >= 0),
    ADD COLUMN repetition_count INTEGER NOT NULL DEFAULT 0 CHECK (repetition_count >= 0),
    ADD COLUMN mastery_score SMALLINT NOT NULL DEFAULT 0 CHECK (mastery_score BETWEEN 0 AND 100);

-- Reconcile vocabulary_reviews with study mode & reaction metrics
ALTER TABLE vocabulary_reviews
    ADD COLUMN study_mode TEXT NOT NULL DEFAULT 'FLASHCARD' CHECK (study_mode IN ('FLASHCARD', 'MULTIPLE_CHOICE', 'TYPE_IN', 'AUDIO_QUIZ')),
    ADD COLUMN response_time_ms INTEGER NOT NULL DEFAULT 0 CHECK (response_time_ms >= 0),
    ADD COLUMN was_correct BOOLEAN NOT NULL DEFAULT TRUE;

-- Create learning_sessions table
CREATE TABLE learning_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language TEXT NOT NULL CHECK (language IN ('JP', 'EN')),
    session_type TEXT NOT NULL CHECK (session_type IN ('REVIEW', 'LEARN_NEW', 'PRACTICE')),
    items_reviewed INTEGER NOT NULL DEFAULT 0 CHECK (items_reviewed >= 0),
    items_correct INTEGER NOT NULL DEFAULT 0 CHECK (items_correct >= 0),
    duration_seconds INTEGER NOT NULL DEFAULT 0 CHECK (duration_seconds >= 0),
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ
);

CREATE INDEX idx_learning_sessions_user ON learning_sessions (user_id, started_at DESC);

-- 4. Journal: Add energy_level, pinned, and word_count
ALTER TABLE journals
    ADD COLUMN energy_level SMALLINT CHECK (energy_level IS NULL OR (energy_level BETWEEN 1 AND 5)),
    ADD COLUMN pinned BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN word_count INTEGER NOT NULL DEFAULT 0 CHECK (word_count >= 0);

-- 5. Performance Indexes
CREATE INDEX IF NOT EXISTS idx_vocabularies_review_queue
    ON vocabularies (user_id, next_review_at ASC);

CREATE INDEX IF NOT EXISTS idx_outbox_pending
    ON outbox_events (created_at) WHERE published_at IS NULL;

-- 6. Cleanup: Drop obsolete legacy asset tables (vehicles & network_nodes)
DROP TABLE IF EXISTS vehicle_logs CASCADE;
DROP TABLE IF EXISTS vehicles CASCADE;
DROP TABLE IF EXISTS network_nodes CASCADE;
