package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (p *Postgres) InsertPlayer(ctx context.Context, plr player.Player) (int, error) {
	var (
		query = `
			INSERT INTO player(
				chat_id,
				created_at,
				payload,
				payload_version,
				payload_updated_at
			) VALUES (
				:chat_id,
				:created_at,
				:payload,
				:payload_version,
				:payload_updated_at
			)
				RETURNING id
		`

		trx    = p.transactor.ExtContext(ctx)
		userID int
	)

	dbPlayer, err := model.ToDBPlayer(plr)
	if err != nil {
		return 0, fmt.Errorf("convert to db player: %w", err)
	}

	query, args, err := sqlx.Named(query, &dbPlayer)
	if err != nil {
		return 0, fmt.Errorf("prepare named: %w", err)
	}

	if err := sqlx.GetContext(
		ctx,
		trx,
		&userID,
		trx.Rebind(query),
		args...,
	); err != nil {
		return 0, fmt.Errorf("get exec: %w", err)
	}

	return userID, nil
}
