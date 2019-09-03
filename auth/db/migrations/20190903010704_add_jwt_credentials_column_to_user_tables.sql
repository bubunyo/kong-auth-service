-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- json is has fast inserter whil jsonp has faster query.
-- Since we dont plan on query the jwt_credentials column,
-- it is safer to user json
ALTER TABLE users ADD COLUMN jwt_credentials json DEFAULT '{}';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users DROP COLUMN IF EXISTS jwt_credentials;
