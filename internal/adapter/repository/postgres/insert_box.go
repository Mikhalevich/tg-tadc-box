package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

func (p *Postgres) InsertBox(
	ctx context.Context,
	domBox box.Box,
) (int, error) {
	var (
		query = `
			INSERT INTO box(
				chat_id,
				status,
				type,
				created_at,
				available_at,
				ready_notification_at,
				completed_at,
				meta
			) VALUES (
				:chat_id,
				:status,
				:type,
				:created_at,
				:available_at,
				:ready_notification_at,
				:completed_at,
				:meta
			)
			RETURNING
				id
		`

		trx = p.transactor.ExtContext(ctx)

		boxID int
	)

	dbBox, err := model.ToDBBox(domBox)
	if err != nil {
		return 0, fmt.Errorf("convert to db box: %w", err)
	}

	query, args, err := sqlx.Named(query, dbBox)
	if err != nil {
		return 0, fmt.Errorf("sqlx named: %w", err)
	}

	if err := sqlx.GetContext(
		ctx,
		trx,
		&boxID,
		trx.Rebind(query),
		args...,
	); err != nil {
		return 0, fmt.Errorf("get context: %w", err)
	}

	return boxID, nil
}
