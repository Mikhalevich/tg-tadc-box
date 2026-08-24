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
	GetRewardByID(
		ctx context.Context,
		id reward.ID,
	) (reward.Reward, error)
	GetRewardsByType(
		ctx context.Context,
		rewardType reward.RewardType,
	) ([]reward.Reward, error)
	GetRewardCountByType(ctx context.Context) (map[reward.RewardType]int, error)
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
