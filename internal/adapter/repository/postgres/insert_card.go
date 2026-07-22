package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
)

// InsertCard insert card to collected_cards table and returns card's count.
func (p *Postgres) InsertCard(
	ctx context.Context,
	receivedCard card.Card,
) (int, error) {
	var (
		query = `
			INSERT INTO collected_cards(
				chat_id,
				reward_id,
				count,
				updated_at
			) VALUES (
				:chat_id,
				:reward_id,
				:count,
				:updated_at
			) ON CONFLICT(chat_id, reward_id)
				DO UPDATE SET
					count = collected_cards.count + EXCLUDED.count,
					updated_at = EXCLUDED.updated_at
			RETURNING
				count
		`

		count int

		trx = p.transactor.ExtContext(ctx)
	)

	query, args, err := sqlx.Named(query, model.ToDBCard(receivedCard))
	if err != nil {
		return 0, fmt.Errorf("sqlx named: %w", err)
	}

	if err := sqlx.GetContext(
		ctx,
		trx,
		&count,
		trx.Rebind(query),
		args...,
	); err != nil {
		return 0, fmt.Errorf("select context: %w", err)
	}

	return count, nil
}
