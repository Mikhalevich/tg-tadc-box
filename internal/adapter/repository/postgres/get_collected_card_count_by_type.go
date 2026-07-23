package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (p *Postgres) GetCollectedCardCountByType(
	ctx context.Context,
	chatID msginfo.ChatID,
) (map[reward.RewardType]int, error) {
	var (
		query = `
			SELECT
				r.type, COUNT(*) AS total
			FROM
				collected_cards AS cc INNER JOIN reward AS r ON cc.reward_id = r.id
			WHERE
				cc.chat_id = $1
			GROUP BY
				r.type
		`

		infos []model.CollectedCardTotalInfo
	)

	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&infos,
		query,
		chatID,
	); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	if len(infos) == 0 {
		return nil, perror.NotFound("cards not found")
	}

	infosMap := make(map[reward.RewardType]int, len(infos))

	for _, info := range infos {
		infosMap[reward.RewardType(info.Type)] = info.Total
	}

	return infosMap, nil
}
