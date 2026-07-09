package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

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
				completed_at
			) VALUES (
				:chat_id,
				:status,
				:type,
				:created_at,
				:available_at,
				:completed_at
			)
			RETURNING
				id
		`

		boxID int
	)

	if err := sqlx.GetContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&boxID,
		query,
	); err != nil {
		return 0, fmt.Errorf("get context: %w", err)
	}

	return boxID, nil
}
