-- name: CreateProject :exec
INSERT INTO projects (
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

-- name: GetProjectByName :one
SELECT *
FROM projects
WHERE name = sqlc.arg(name);

-- name: ListProjects :many
SELECT *
FROM projects
ORDER BY name ASC;