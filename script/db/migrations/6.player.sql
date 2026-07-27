-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE player(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    payload JSONB NOT NULL,
    payload_version INTEGER NOT NULL DEFAULT 0,
    payload_updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX player_chat_id_idx ON player(chat_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE player;
