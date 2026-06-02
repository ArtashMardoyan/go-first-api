-- +goose Up
CREATE TABLE posts (
    "id"        TEXT        NOT NULL PRIMARY KEY,
    "title"     TEXT        NOT NULL,
    "body"      TEXT        NOT NULL DEFAULT '',
    "status"    TEXT        NOT NULL DEFAULT 'unpublished',
    "userId"    TEXT        NOT NULL REFERENCES users ("id") ON DELETE CASCADE,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX posts_userId_idx ON posts ("userId");

-- +goose Down
DROP TABLE posts;