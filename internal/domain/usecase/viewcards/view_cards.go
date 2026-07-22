package viewcards

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Repository interface {
	GetCollectedCardAfterID(ctx context.Context, afterID card.ID) (card.Card, error)
}

type RewardProvider interface {
	GetRewardByID(ctx context.Context, id reward.ID) (reward.Reward, error)
}

type Notifier interface {
	NoCollectedCards(ctx context.Context, chatID msginfo.ChatID) error
	ShowCollectedReward(
		ctx context.Context,
		chatID msginfo.ChatID,
		rew reward.Reward,
		count int,
	) error
}

type ViewCards struct {
	repo           Repository
	rewardProvider RewardProvider
	notifier       Notifier
}

func New(
	repo Repository,
	rewardProvider RewardProvider,
	notifier Notifier,
) *ViewCards {
	return &ViewCards{
		repo:           repo,
		rewardProvider: rewardProvider,
		notifier:       notifier,
	}
}
