-- migrate:up

CREATE TABLE authors
(
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    bio  TEXT NOT NULL
);

CREATE INDEX ON authors (id);

-- migrate:down
DROP TABLE IF EXISTS authors;
