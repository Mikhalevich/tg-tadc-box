package rewardsvc

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type RewardPercent struct {
	Rare      int
	Epic      int
	Legendary int
}

type Repository interface {
	GetRewardsByType(
		ctx context.Context,
		rewardType reward.RewardType,
	) ([]reward.Reward, error)
}

type Service struct {
	repo          Repository
	rewardPercent map[box.Type]RewardPercent
}

func New(
	rewardPercent map[box.Type]RewardPercent,
	repo Repository,
) *Service {
	return &Service{
		rewardPercent: rewardPercent,
		repo:          repo,
	}
}
