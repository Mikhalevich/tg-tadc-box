package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (p *Postgres) GetCollectedCardByPos(
	ctx context.Context,
	chatID msginfo.ChatID,
	rewardType reward.RewardType,
	position int,
) (card.Card, error) {
	var (
		query = `
			WITH cards_by_pos AS (
				SELECT
					cc.id,
					cc.chat_id,
					cc.reward_id,
					cc.count,
					cc.updated_at,
					ROW_NUMBER() OVER (ORDER BY cc.id) AS position
				FROM
					collected_cards AS cc INNER JOIN reward AS r ON cc.reward_id = r.id
				WHERE
					cc.chat_id = $1 AND
					r.type = $2
					
			)
			SELECT
				id,
				chat_id,
				reward_id,
				count,
				updated_at
			FROM
				cards_by_pos
			WHERE
				position = $3
		`

		dbCard model.Card
	)

	if err := sqlx.GetContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&dbCard,
		query,
		chatID.Int64(),
		rewardType.String(),
		position,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return card.Card{}, perror.NotFound("cards not found")
		}

		return card.Card{}, fmt.Errorf("get context: %w", err)
	}

	return dbCard.ToDomCard(), nil
}
