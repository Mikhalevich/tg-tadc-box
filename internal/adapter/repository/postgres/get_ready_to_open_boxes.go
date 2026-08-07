package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

func (p *Postgres) GetReadyToOpenBoxes(
	ctx context.Context,
	now time.Time,
	limit int,
) ([]box.Box, error) {
	var (
		query = `
			SELECT
				id,
				chat_id,
				status,
				type,
				created_at,
				available_at,
				ready_notification_at,
				completed_at,
				meta
			FROM
				box
			WHERE
				status = $1 AND
				available_at <= $2 AND
				ready_notification_at IS NULL
			LIMIT
				$3
			FOR UPDATE SKIP LOCKED
		`

		dbBoxes []model.Box
	)

	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&dbBoxes,
		query,
		box.StatusInProgress,
		now,
		limit,
	); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	domBoxes, err := model.ToDomBoxes(dbBoxes)
	if err != nil {
		return nil, fmt.Errorf("convert to dom boxes: %w", err)
	}

	return domBoxes, nil
}
