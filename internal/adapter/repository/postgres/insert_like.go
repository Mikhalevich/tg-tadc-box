package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/like"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (p *Postgres) InsertLike(
	ctx context.Context,
	rwdLike like.Like,
) error {
	var (
		query = `
			INSERT INTO likes(
				box_id,
				chat_id,
				reward_id,
				type,
				created_at
			) VALUES (
				:box_id,
				:chat_id,
				:reward_id,
				:type,
				:created_at
			)
		`
	)

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		model.ToDBLike(rwdLike),
	)

	if err != nil {
		return fmt.Errorf("exec context: %w", err)
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
