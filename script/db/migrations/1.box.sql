-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE box_status AS ENUM (
    'pending',
    'in_progress',
    'opened',
    'canceled'
);

CREATE TYPE box_type AS ENUM (
    'normal'
);

CREATE TABLE box(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL,
    status box_status NOT NULL,
    type box_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    available_at TIMESTAMPTZ NOT NULL,
    ready_notification_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

CREATE INDEX box_active_boxes ON box(chat_id, status) WHERE status IN ('pending', 'in_progress');
CREATE INDEX box_ready_to_open_notification ON box(status, available_at) WHERE status IN ('pending', 'in_progress') AND ready_notification_at IS NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE box;
DROP TYPE box_type;
DROP TYPE box_status;
