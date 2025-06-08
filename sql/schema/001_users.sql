-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    body STRING NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE
);



-- +goose Down
DROP TABLE users;
DROP TABLE chirps;

