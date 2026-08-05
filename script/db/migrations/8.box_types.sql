-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TYPE box_type RENAME VALUE 'normal' TO 'common';
ALTER TYPE box_type ADD VALUE 'rare';
ALTER TYPE box_type ADD VALUE 'epic';
ALTER TYPE box_type ADD VALUE 'legendary';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
