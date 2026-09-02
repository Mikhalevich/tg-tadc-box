-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE referral(
    chat_id BIGINT PRIMARY KEY,
    invited_by_chat_id BIGINT,
    code TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX referral_code_idx ON referral(code);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE referral;
