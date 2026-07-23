package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (p *Postgres) GetRewardByID(
	ctx context.Context,
	rewardID reward.ID,
) (reward.Reward, error) {
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
				id = $1
		`

		dbReward model.Reward
	)

	if err := sqlx.GetContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&dbReward,
		query,
		rewardID.Int(),
	); err != nil {
		return reward.Reward{}, fmt.Errorf("get context: %w", err)
	}

	return dbReward.ToDom(), nil
}
