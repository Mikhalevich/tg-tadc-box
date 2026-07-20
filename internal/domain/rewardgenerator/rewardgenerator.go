package rewardgenerator

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

const (
	legendaryPercent = 5
	epicPercent      = 10
	rarePercent      = 30
)

type RewardsGetter interface {
	GetRewardsByType(
		ctx context.Context,
		rewardType reward.RewardType,
	) ([]reward.Reward, error)
}

type RewardGenerator struct {
	rewardsGetter RewardsGetter
}

func New(
	rewardsGetter RewardsGetter,
) RewardGenerator {
	return RewardGenerator{
		rewardsGetter: rewardsGetter,
	}
}

func (r RewardGenerator) Generate(ctx context.Context) (reward.Reward, error) {
	rewardType := pickRewardType()

	rewards, err := r.rewardsGetter.GetRewardsByType(ctx, rewardType)
	if err != nil {
		return reward.Reward{}, fmt.Errorf("get all rewards: %w", err)
	}

	if len(rewards) == 0 {
		return reward.Reward{}, fmt.Errorf("no rewards by type %q", rewardType.String())
	}

	return pickReward(rewards), nil
}

func percent() int {
	//nolint:gosec,mnd
	return rand.IntN(100) + 1
}

func pickRewardType() reward.RewardType {
	roll := percent()

	switch {
	case roll <= legendaryPercent:
		return reward.RewardTypeLegendary

	case roll <= epicPercent:
		return reward.RewardTypeEpic

	case roll <= rarePercent:
		return reward.RewardTypeRare
	}

	return reward.RewardTypeCommon
}

func pickReward(
	rewards []reward.Reward,
) reward.Reward {
	//nolint:gosec
	return rewards[rand.IntN(len(rewards))]
}
