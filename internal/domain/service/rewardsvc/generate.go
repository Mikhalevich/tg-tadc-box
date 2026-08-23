package rewardsvc

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (s *Service) Generate(
	ctx context.Context,
	boxType box.Type,
) (reward.Reward, error) {
	rewardPercent, ok := s.rewardPercent[boxType]
	if !ok {
		return reward.Reward{}, fmt.Errorf("no such box reward %q", boxType.String())
	}

	rewardType := pickRewardType(rewardPercent)

	rewards, err := s.repo.GetRewardsByType(ctx, rewardType)
	if err != nil {
		return reward.Reward{}, fmt.Errorf("get rewards by type: %w", err)
	}

	if len(rewards) == 0 {
		return reward.Reward{}, fmt.Errorf("no rewards by type %q", rewardType.String())
	}

	return pickRandomReward(rewards), nil
}

func percent() int {
	//nolint:gosec,mnd
	return rand.IntN(100) + 1
}

func pickRewardType(rewardPercent RewardPercent) reward.RewardType {
	roll := percent()

	switch {
	case roll <= rewardPercent.Legendary:
		return reward.RewardTypeLegendary

	case roll <= rewardPercent.Epic:
		return reward.RewardTypeEpic

	case roll <= rewardPercent.Rare:
		return reward.RewardTypeRare
	}

	return reward.RewardTypeCommon
}

func pickRandomReward(
	rewards []reward.Reward,
) reward.Reward {
	//nolint:gosec
	return rewards[rand.IntN(len(rewards))]
}
