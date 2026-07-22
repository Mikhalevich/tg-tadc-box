-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE reward_type AS ENUM (
    'common',
    'rare',
    'epic',
    'legendary'
);

CREATE TABLE reward(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    type reward_type NOT NULL,
    uri TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX reward_uri_idx ON reward(uri);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE reward;
DROP TYPE reward_type;
