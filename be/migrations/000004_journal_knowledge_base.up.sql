ALTER TABLE journals
    ADD COLUMN mood TEXT,
    ADD COLUMN published_date DATE NOT NULL DEFAULT CURRENT_DATE,
    ADD CONSTRAINT journals_mood_check CHECK (mood IS NULL OR mood IN ('VERY_NEGATIVE', 'NEGATIVE', 'NEUTRAL', 'POSITIVE', 'VERY_POSITIVE'));

CREATE INDEX journals_user_published_idx ON journals (user_id, published_date DESC, id DESC);

CREATE TABLE journal_links (
    journal_id UUID NOT NULL REFERENCES journals(id) ON DELETE CASCADE,
    linked_id UUID NOT NULL REFERENCES journals(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (journal_id, linked_id),
    CHECK (journal_id <> linked_id)
);

CREATE INDEX journal_links_user_linked_idx ON journal_links (user_id, linked_id);

CREATE OR REPLACE FUNCTION journals_tsvector_trigger() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.content, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER journals_tsvector_update
    BEFORE INSERT OR UPDATE OF title, content ON journals
    FOR EACH ROW EXECUTE FUNCTION journals_tsvector_trigger();

UPDATE journals SET title = title;

ALTER TABLE journals ALTER COLUMN published_date DROP DEFAULT;
