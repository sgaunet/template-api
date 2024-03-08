
-- name: CreateBook :one
INSERT INTO books (title, author_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetBook :one
SELECT *
FROM books
WHERE id = $1
LIMIT 1;

-- name: UpdateTitleBook :one
UPDATE books
SET title = $2
WHERE id = $1
RETURNING *;

-- name: DeleteBook :exec
DELETE
FROM books
WHERE id = $1;

-- name: ListBooks :many
SELECT *
FROM books
ORDER BY title;
