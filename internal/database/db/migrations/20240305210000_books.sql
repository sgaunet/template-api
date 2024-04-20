-- migrate:up

CREATE TABLE books
(
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(32) NOT NULL,
    author_id   BIGINT NOT NULL,
    FOREIGN KEY (author_id) REFERENCES authors (id)
);

CREATE INDEX ON books (id);

-- migrate:down
DROP TABLE IF EXISTS books;