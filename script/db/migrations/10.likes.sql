-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE like_type AS ENUM (
    'like',
    'dislike'
);

CREATE TABLE likes(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    box_id INTEGER NOT NULL,
    chat_id BIGINT NOT NULL,
    reward_id INTEGER NOT NULL,
    type like_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT likes_reward_id_fk FOREIGN KEY(reward_id) REFERENCES reward(id)

);

CREATE UNIQUE INDEX likes_box_id_idx ON likes(box_id);
CREATE INDEX likes_reward_id_like_type_idx ON likes(reward_id, type);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE likes;
DROP TYPE like_type;
