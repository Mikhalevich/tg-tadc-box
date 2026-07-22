package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (p *Postgres) GetCollectedCardAfterID(
	ctx context.Context,
	afterID card.ID,
) (card.Card, error) {
	var (
		query = `
			SELECT
				id,
				chat_id,
				reward_id,
				count,
				updated_at
			FROM
				collected_cards
			WHERE
				id > $1
			ORDER BY
				id
			LIMIT
				1
		`

		dbCard model.Card
	)

	if err := sqlx.GetContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&dbCard,
		query,
		afterID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return card.Card{}, perror.NotFound("cards not found")
		}

		return card.Card{}, fmt.Errorf("get context: %w", err)
	}

	return dbCard.ToDomCard(), nil
}
