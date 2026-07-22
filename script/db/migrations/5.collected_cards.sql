-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE collected_cards(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL,
    reward_id INTEGER NOT NULL,
    count INTEGER NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT collected_cards_reward_id_fk FOREIGN KEY(reward_id) REFERENCES reward(id)
);

CREATE UNIQUE INDEX collected_cards_chat_id_reward_id_idx ON collected_cards(chat_id, reward_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE collected_cards;
