-- name: CreateEntry :exec
INSERT INTO entries (
    id,
    description,
    project_id,
    started_at,
    stopped_at,
    created_at,
    updated_at
) VALUES (
    sqlc.arg(id),
    sqlc.arg(description),
    sqlc.narg(project_id),
    sqlc.arg(started_at),
    sqlc.narg(stopped_at),
    sqlc.arg(created_at),
    sqlc.arg(updated_at)
);

-- name: CreateEntryTag :exec
INSERT INTO entry_tags (
    entry_id,
    tag_id,
    position
) VALUES (
    sqlc.arg(entry_id),
    sqlc.arg(tag_id),
    sqlc.arg(position)
);

-- name: GetActiveEntry :one
SELECT
    id,
    description,
    project_id,
    started_at,
    stopped_at,
    created_at,
    updated_at
FROM entries
WHERE stopped_at IS NULL
LIMIT 1;

-- name: StopActiveEntry :one
UPDATE entries
SET
    stopped_at = sqlc.arg(stopped_at),
    updated_at = sqlc.arg(updated_at)
WHERE stopped_at IS NULL
RETURNING
    id,
    description,
    project_id,
    started_at,
    stopped_at,
    created_at,
    updated_at;

-- name: ListEntriesByStartRange :many
SELECT
    id,
    description,
    project_id,
    started_at,
    stopped_at,
    created_at,
    updated_at
FROM entries
WHERE started_at >= sqlc.arg(started_at_from)
    AND started_at < sqlc.arg(started_at_to)
ORDER BY started_at ASC, created_at ASC;

-- name: ListRecentEntries :many
SELECT
    id,
    description,
    project_id,
    started_at,
    stopped_at,
    created_at,
    updated_at
FROM entries
ORDER BY started_at DESC, created_at DESC
LIMIT sqlc.arg(limit);

-- name: ListTagsForEntry :many
SELECT tags.name
FROM entry_tags
JOIN tags ON tags.id = entry_tags.tag_id
WHERE entry_tags.entry_id = sqlc.arg(entry_id)
ORDER BY entry_tags.position ASC;
