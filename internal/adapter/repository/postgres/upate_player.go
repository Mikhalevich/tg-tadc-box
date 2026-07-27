package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (p *Postgres) UpdatePlayer(ctx context.Context, plr player.Player) error {
	var (
		query = `
			UPDATE player SET
				payload = :payload,
				payload_version = payload_version + 1,
				payload_updated_at = :payload_updated_at
			WHERE
				id = :id AND
				payload_version = :payload_version
		`
	)

	dbPlayer, err := model.ToDBPlayer(plr)
	if err != nil {
		return fmt.Errorf("convert to db user: %w", err)
	}

	res, err := sqlx.NamedExecContext(
		ctx,
		p.transactor.ExtContext(ctx),
		query,
		&dbPlayer,
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
