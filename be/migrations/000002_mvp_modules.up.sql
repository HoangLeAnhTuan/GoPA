CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL CHECK (char_length(trim(title)) BETWEEN 1 AND 200),
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL CHECK (status IN ('TODO', 'IN_PROGRESS', 'DONE')),
    priority TEXT NOT NULL CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH')),
    category TEXT NOT NULL CHECK (category IN ('WORK', 'STUDY', 'LIFE')),
    due_date DATE,
    position BIGINT NOT NULL CHECK (position >= 0),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX tasks_user_status_position_idx ON tasks (user_id, status, position, created_at DESC);

CREATE TABLE vocabularies (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language TEXT NOT NULL CHECK (language IN ('JP', 'EN')),
    word TEXT NOT NULL CHECK (char_length(trim(word)) BETWEEN 1 AND 300),
    reading TEXT,
    meaning TEXT NOT NULL CHECK (char_length(trim(meaning)) BETWEEN 1 AND 1000),
    example_sentence TEXT,
    current_box SMALLINT NOT NULL CHECK (current_box BETWEEN 1 AND 5),
    next_review_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (user_id, language, word)
);

CREATE INDEX vocabularies_review_queue_idx ON vocabularies (user_id, next_review_at);

CREATE TABLE vocabulary_reviews (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL UNIQUE,
    vocabulary_id UUID NOT NULL REFERENCES vocabularies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quality SMALLINT NOT NULL CHECK (quality BETWEEN 0 AND 5),
    reviewed_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX vocabulary_reviews_vocabulary_idx ON vocabulary_reviews (vocabulary_id, reviewed_at DESC);

CREATE TABLE vehicles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    make TEXT NOT NULL CHECK (char_length(trim(make)) BETWEEN 1 AND 100),
    model TEXT NOT NULL CHECK (char_length(trim(model)) BETWEEN 1 AND 100),
    year SMALLINT NOT NULL CHECK (year BETWEEN 1886 AND 2100),
    nickname TEXT,
    current_mileage INTEGER NOT NULL CHECK (current_mileage >= 0),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX vehicles_user_idx ON vehicles (user_id, created_at DESC);

CREATE TABLE vehicle_logs (
    id UUID PRIMARY KEY,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_type TEXT NOT NULL CHECK (log_type IN ('MAINTENANCE', 'FUEL', 'UPGRADE')),
    mileage INTEGER NOT NULL CHECK (mileage >= 0),
    description TEXT NOT NULL CHECK (char_length(trim(description)) BETWEEN 1 AND 2000),
    cost NUMERIC(14, 2) NOT NULL CHECK (cost >= 0),
    log_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX vehicle_logs_vehicle_date_idx ON vehicle_logs (vehicle_id, log_date DESC, created_at DESC);

CREATE TABLE network_nodes (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_name TEXT NOT NULL CHECK (char_length(trim(device_name)) BETWEEN 1 AND 200),
    mac_address TEXT,
    static_ip TEXT,
    vlan_tag SMALLINT CHECK (vlan_tag BETWEEN 1 AND 4094),
    connection_type TEXT NOT NULL CHECK (connection_type IN ('LAN', 'WIFI')),
    parent_node_id UUID REFERENCES network_nodes(id) ON DELETE SET NULL,
    location TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CHECK (parent_node_id IS NULL OR parent_node_id <> id)
);

CREATE INDEX network_nodes_user_parent_idx ON network_nodes (user_id, parent_node_id, created_at);

CREATE TABLE pomodoro_history (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ NOT NULL,
    duration_seconds INTEGER NOT NULL CHECK (duration_seconds > 0),
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX pomodoro_history_user_started_idx ON pomodoro_history (user_id, started_at DESC);

CREATE TABLE journals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL CHECK (char_length(trim(title)) BETWEEN 1 AND 300),
    content TEXT NOT NULL,
    tags JSONB NOT NULL DEFAULT '[]'::jsonb CHECK (jsonb_typeof(tags) = 'array'),
    search_vector TSVECTOR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX journals_user_updated_idx ON journals (user_id, updated_at DESC);
CREATE INDEX journals_search_vector_idx ON journals USING GIN (search_vector);
CREATE INDEX journals_tags_idx ON journals USING GIN (tags);

CREATE TABLE processed_events (
    event_id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    routing_key TEXT NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX outbox_events_unpublished_idx ON outbox_events (created_at) WHERE published_at IS NULL;
