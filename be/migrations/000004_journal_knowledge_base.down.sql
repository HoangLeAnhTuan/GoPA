DROP TRIGGER IF EXISTS journals_tsvector_update ON journals;
DROP FUNCTION IF EXISTS journals_tsvector_trigger();
DROP TABLE IF EXISTS journal_links;
DROP INDEX IF EXISTS journals_user_published_idx;
ALTER TABLE journals DROP CONSTRAINT IF EXISTS journals_mood_check;
ALTER TABLE journals DROP COLUMN IF EXISTS published_date;
ALTER TABLE journals DROP COLUMN IF EXISTS mood;
