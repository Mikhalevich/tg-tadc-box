package rewardgenerator

import (
	"fmt"
	"math/rand/v2"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

const (
	legendaryPercent = 5
	epicPercent      = 10
	rarePercent      = 30
)

type RewardGenerator struct {
	rewardIDs map[reward.RewardType][]reward.ID
}

func New(
	ids map[reward.RewardType][]reward.ID,
) RewardGenerator {
	return RewardGenerator{
		rewardIDs: ids,
	}
}

func (r RewardGenerator) Generate() (reward.Reward, error) {
	rewardType := generateRewardType()

	rewardID, err := r.generateRewardID(rewardType)
	if err != nil {
		return reward.Reward{}, fmt.Errorf("generate reward id: %w", err)
	}

	return reward.Reward{
		ID:   rewardID,
		Type: rewardType,
	}, nil
}

func percent() int {
	//nolint:gosec,mnd
	return rand.IntN(100) + 1
}

func generateRewardType() reward.RewardType {
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

func (r RewardGenerator) generateRewardID(rewardType reward.RewardType) (reward.ID, error) {
	ids, ok := r.rewardIDs[rewardType]
	if !ok {
		return 0, fmt.Errorf("invalid reward type: %s", rewardType)
	}

	//nolint:gosec
	return ids[rand.IntN(len(ids))], nil
}
