package rewardsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (s *Service) GetRewardCountByType(
	ctx context.Context,
) (map[reward.RewardType]int, error) {
	counts, err := s.repo.GetRewardCountByType(ctx)
	if err != nil {
		return nil, fmt.Errorf("get reward count by type: %w", err)
	}

	return counts, nil
}
