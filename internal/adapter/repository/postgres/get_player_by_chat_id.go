package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (p *Postgres) GetPlayerByChatID(
	ctx context.Context,
	chatID msginfo.ChatID,
) (player.Player, error) {
	var (
		query = `
			SELECT
				id,
				chat_id,
				created_at,
				payload,
				payload_version,
				payload_updated_at
			FROM
				player
			WHERE
				chat_id = $1
			FOR UPDATE
		`

		dbPlayer model.Player
	)

	if err := sqlx.GetContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&dbPlayer,
		query,
		chatID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return player.Player{}, perror.NotFound("player not found")
		}

		return player.Player{}, fmt.Errorf("get context: %w", err)
	}

	domPlayer, err := dbPlayer.ToDom()
	if err != nil {
		return player.Player{}, fmt.Errorf("convert to dom player: %w", err)
	}

	return domPlayer, nil
}
