package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (p *Postgres) GetBoxByID(
	ctx context.Context,
	boxID int,
) (box.Box, error) {
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
				completed_at
			FROM
				box
			WHERE
				id = $1
			FOR UPDATE
		`

		dbBox model.Box
	)

	if err := sqlx.GetContext(ctx, p.transactor.ExtContext(ctx), &dbBox, query, boxID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return box.Box{}, perror.NotFound("box not found")
		}

		return box.Box{}, fmt.Errorf("get context: %w", err)
	}

	return dbBox.ToDomBox(), nil
}
