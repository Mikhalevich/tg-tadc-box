-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE button ADD COLUMN style TEXT NOT NULL DEFAULT '';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE button DROP COLUMN style;
