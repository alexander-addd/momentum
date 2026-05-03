-- +goose Up
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    UNIQUE (name)
);

CREATE TABLE tags (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    UNIQUE (name)
);

CREATE TABLE entries (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL CHECK (length(trim(description)) > 0),
    project_id TEXT,
    started_at INTEGER NOT NULL,
    stopped_at INTEGER,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    CHECK (stopped_at IS NULL OR stopped_at >= started_at),
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX entries_one_active_idx
    ON entries ((1))
    WHERE stopped_at IS NULL;

CREATE INDEX entries_started_at_idx
    ON entries (started_at DESC);

CREATE TABLE entry_tags (
    entry_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    position INTEGER NOT NULL CHECK (position >= 0),
    PRIMARY KEY (entry_id, tag_id),
    UNIQUE (entry_id, position),
    FOREIGN KEY (entry_id) REFERENCES entries (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
);

CREATE INDEX entry_tags_tag_idx
    ON entry_tags (tag_id);

-- +goose Down
DROP INDEX entry_tags_tag_idx;
DROP TABLE entry_tags;
DROP INDEX entries_started_at_idx;
DROP INDEX entries_one_active_idx;
DROP TABLE entries;
DROP TABLE tags;
DROP TABLE projects;
