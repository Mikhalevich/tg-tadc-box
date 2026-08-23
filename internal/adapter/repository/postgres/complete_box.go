package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (p *Postgres) CompleteBox(
	ctx context.Context,
	boxID box.ID,
	status box.Status,
	completedAt time.Time,
) error {
	var (
		query = `
			UPDATE box SET
				status = :status,
				completed_at = :completed_at
			WHERE
				id = :id AND
				status <> :status
		`
	)

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		map[string]any{
			"status":       status.String(),
			"completed_at": completedAt,
		},
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
