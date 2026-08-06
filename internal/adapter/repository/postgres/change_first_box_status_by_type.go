package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (p *Postgres) ChangeFirstBoxStatusByType(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	newStatus box.Status,
	previousStatus box.Status,
	availableAt time.Time,
) error {
	var (
		query = `
			UPDATE box SET
				status = :status,
				available_at = :available_at
			WHERE
				id = (
					SELECT
						id
					FROM
						box
					WHERE
						chat_id = :chat_id AND
						type = :type AND
						status = :previous_status
					ORDER BY
						created_at
					LIMIT
						1
				)
		`
	)

	if _, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		map[string]any{
			"chat_id":         chatID.Int64(),
			"type":            boxType.String(),
			"status":          newStatus.String(),
			"previous_status": previousStatus.String(),
			"available_at":    availableAt,
		},
	); err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	return nil
}
