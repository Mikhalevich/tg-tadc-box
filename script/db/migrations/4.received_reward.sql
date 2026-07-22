-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE received_reward(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL,
    reward_id INTEGER NOT NULL,
    box_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT received_reward_reward_id_fk FOREIGN KEY(reward_id) REFERENCES reward(id),
    CONSTRAINT received_reward_box_id_fk FOREIGN KEY(box_id) REFERENCES box(id)
);

CREATE INDEX received_reward_reward_id_idx ON received_reward(reward_id);
CREATE INDEX received_reward_box_id_idx ON received_reward(box_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE received_reward;
