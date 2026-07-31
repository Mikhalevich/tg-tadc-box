package rewardgenerator

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type RewardPercent struct {
	Rare      int
	Epic      int
	Legendary int
}

type RewardsGetter interface {
	GetRewardsByType(
		ctx context.Context,
		rewardType reward.RewardType,
	) ([]reward.Reward, error)
}

type RewardGenerator struct {
	rewardsGetter RewardsGetter
	rewardPercent map[box.Type]RewardPercent
}

func New(
	rewardsGetter RewardsGetter,
	rewardPercent map[box.Type]RewardPercent,
) RewardGenerator {
	return RewardGenerator{
		rewardsGetter: rewardsGetter,
		rewardPercent: rewardPercent,
	}
}

func (r RewardGenerator) Generate(
	ctx context.Context,
	boxType box.Type,
) (reward.Reward, error) {
	rewardPercent, ok := r.rewardPercent[boxType]
	if !ok {
		return reward.Reward{}, fmt.Errorf("no such box reward %q", boxType.String())
	}

	rewardType := pickRewardType(rewardPercent)

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

func pickReward(
	rewards []reward.Reward,
) reward.Reward {
	//nolint:gosec
	return rewards[rand.IntN(len(rewards))]
}
