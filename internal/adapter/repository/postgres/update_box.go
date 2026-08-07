package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (p *Postgres) UpdateBox(
	ctx context.Context,
	domBox box.Box,
) error {
	var (
		query = `
			UPDATE box SET
				status = :status,
				available_at = :available_at,
				completed_at = :completed_at
			WHERE
				id = :id
		`
	)

	dbBox, err := model.ToDBBox(domBox)
	if err != nil {
		return fmt.Errorf("convert to db box: %w", err)
	}

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		dbBox,
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
