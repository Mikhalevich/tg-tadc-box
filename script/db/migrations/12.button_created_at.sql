-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE button ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE button DROP COLUMN created_at;
