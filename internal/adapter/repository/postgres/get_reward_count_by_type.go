package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (p *Postgres) GetRewardCountByType(
	ctx context.Context,
) (map[reward.RewardType]int, error) {
	var (
		query = `
			SELECT
				type, COUNT(*) AS total
			FROM
				reward
			GROUP BY
				type
		`

		infos []model.RewardTotalInfo
	)

	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&infos,
		query,
	); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	if len(infos) == 0 {
		return nil, perror.NotFound("rewards not found")
	}

	infosMap := make(map[reward.RewardType]int, len(infos))

	for _, info := range infos {
		infosMap[reward.RewardType(info.Type)] = info.Total
	}

	return infosMap, nil
}
