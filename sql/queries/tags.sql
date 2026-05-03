-- name: CreateTag :exec
INSERT INTO tags (
    id,
    name,
    created_at,
    updated_at
) VALUES (
    sqlc.arg(id),
    sqlc.arg(name),
    sqlc.arg(created_at),
    sqlc.arg(updated_at)
);

-- name: GetTagByName :one
SELECT *
FROM tags
WHERE name = sqlc.arg(name);

-- name: ListTags :many
SELECT *
FROM tags
ORDER BY name ASC;