-- +gooseUp

ALTER TABLE users
ADD COLUMN is_chirpy_red BOOL DEFAULT false;
