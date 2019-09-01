-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE users (id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(), email text UNIQUE, password text);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users;
