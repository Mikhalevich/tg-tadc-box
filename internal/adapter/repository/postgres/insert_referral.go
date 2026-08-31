package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

func (p *Postgres) InsertReferral(
	ctx context.Context,
	ref referral.Referral,
) error {
	var (
		query = `
			INSERT INTO referral(
				chat_id,
				invited_by_chat_id,
				code,
				created_at
			) VALUES (
				:chat_id,
				:invited_by_chat_id,
				:code,
				:created_at
			)
		`
	)

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		model.ToDBReferral(ref),
	)

	if err != nil {
		return fmt.Errorf("named exec: %w", err)
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
