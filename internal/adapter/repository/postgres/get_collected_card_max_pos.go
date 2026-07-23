package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (p *Postgres) GetCollectedCardMaxPos(
	ctx context.Context,
	chatID msginfo.ChatID,
	rewardType reward.RewardType,
) (int, error) {
	var (
		query = `
			WITH cards_by_pos AS (
				SELECT
					ROW_NUMBER() OVER (ORDER BY cc.id) AS position
				FROM
					collected_cards AS cc INNER JOIN reward AS r ON cc.reward_id = r.id
				WHERE
					cc.chat_id = $1 AND
					r.type = $2
			)
			SELECT
				COALESCE(MAX(position), 0)
			FROM
				cards_by_pos
		`

		maxPos int
	)

	if err := sqlx.GetContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&maxPos,
		query,
		chatID.Int64(),
		rewardType,
	); err != nil {
		return 0, fmt.Errorf("get context: %w", err)
	}

	return maxPos, nil
}
