package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (p *Postgres) InsertReceivedReward(
	ctx context.Context,
	rwd reward.ReceivedReward,
) error {
	var (
		query = `
			INSERT INTO received_reward(
				chat_id,
				reward_id,
				box_id,
				created_at
			) VALUES (
				:chat_id,
				:reward_id,
				:box_id,
				:created_at
			)
		`
	)

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		model.ToDBReceivedReward(rwd),
	)
	if err != nil {
		return fmt.Errorf("named exec context: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return perror.NoRowsUpdated()
	}

	return nil
}
