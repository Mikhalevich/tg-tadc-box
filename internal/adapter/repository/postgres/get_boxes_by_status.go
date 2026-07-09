package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (p *Postgres) GetBoxesByStatus(
	ctx context.Context,
	chatID msginfo.ChatID,
	statuses ...box.Status,
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
				completed_at
			FROM
				box
			WHERE
				chat_id = ? AND
				status IN (?)
		`

		trx = p.transactor.ExtContext(ctx)

		boxes []model.Box
	)

	query, args, err := sqlx.In(query, chatID, statuses)
	if err != nil {
		return nil, fmt.Errorf("sqlx in: %w", err)
	}

	if err := sqlx.SelectContext(
		ctx,
		trx,
		&boxes,
		trx.Rebind(query),
		args...,
	); err != nil {
		return nil, fmt.Errorf("get context: %w", err)
	}

	return model.ToDomBoxes(boxes), nil
}
