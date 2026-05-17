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
    e.id,
    e.description,
    e.project_id,
    p.name AS project_name,
    e.started_at,
    e.stopped_at,
    e.created_at,
    e.updated_at
FROM entries e
LEFT JOIN projects p ON p.id = e.project_id
WHERE stopped_at IS NULL
LIMIT 1;

-- name: GetEntryByID :one
SELECT
    e.id,
    e.description,
    e.project_id,
    p.name AS project_name,
    e.started_at,
    e.stopped_at,
    e.created_at,
    e.updated_at
FROM entries e
LEFT JOIN projects p ON p.id = e.project_id
WHERE e.id = sqlc.arg(id)
LIMIT 1;

-- name: GetEntries :many
SELECT
    e.id,
    e.description,
    e.project_id,
    p.name AS project_name,
    e.started_at,
    e.stopped_at,
    e.created_at,
    e.updated_at
FROM entries e
LEFT JOIN projects p ON p.id = e.project_id
ORDER BY e.created_at DESC
LIMIT sqlc.arg(limit);

-- name: StopActiveEntry :execrows
UPDATE entries
SET
    stopped_at = sqlc.arg(stopped_at),
    updated_at = sqlc.arg(updated_at)
WHERE stopped_at IS NULL
    AND id = sqlc.arg(id);

