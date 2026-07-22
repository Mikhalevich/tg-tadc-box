package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (p *Postgres) GetRewardsByType(
	ctx context.Context,
	rewardType reward.RewardType,
) ([]reward.Reward, error) {
	var (
		query = `
			SELECT
				id,
				type,
				uri,
				created_at
			FROM
				reward
			WHERE
				type = $1
		`

		dbRewards []model.Reward
	)

	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&dbRewards,
		query,
		rewardType.String(),
	); err != nil {
		return nil, fmt.Errorf("select context: %w", err)
	}

	return model.ToDomRewards(dbRewards), nil
}
