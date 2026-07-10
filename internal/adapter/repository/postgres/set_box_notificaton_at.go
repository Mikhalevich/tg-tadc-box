package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (p *Postgres) SetBoxReadyNotificationAt(
	ctx context.Context,
	ids []int,
	notificationAt time.Time,
) error {
	var (
		query = `
			UPDATE box SET
				ready_notification_at = ?
			WHERE
				id iN (?)
		`

		trx = p.transactor.ExtContext(ctx)
	)

	query, args, err := sqlx.In(query, notificationAt, ids)
	if err != nil {
		return fmt.Errorf("sqlx in: %w", err)
	}

	res, err := trx.ExecContext(ctx, trx.Rebind(query), args...)
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
