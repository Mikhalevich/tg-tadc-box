package rewardsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (s *Service) GetRewardByID(
	ctx context.Context,
	id reward.ID,
) (reward.Reward, error) {
	rwd, err := s.repo.GetRewardByID(ctx, id)
	if err != nil {
		return reward.Reward{}, fmt.Errorf("get reward by id: %w", err)
	}

	return rwd, nil
}
