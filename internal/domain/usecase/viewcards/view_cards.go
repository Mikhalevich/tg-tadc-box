package viewcards

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Repository interface {
	GetCollectedCardCountByType(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (map[reward.RewardType]int, error)
	GetCollectedCardByPos(
		ctx context.Context,
		chatID msginfo.ChatID,
		rewardType reward.RewardType,
		pos int,
	) (card.Card, error)
	GetCollectedCardMaxPos(
		ctx context.Context,
		chatID msginfo.ChatID,
		rewardType reward.RewardType,
	) (int, error)
}

type RewardProvider interface {
	GetRewardByID(ctx context.Context, id reward.ID) (reward.Reward, error)
	GetRewardCountByType(ctx context.Context) (map[reward.RewardType]int, error)
}

type Notifier interface {
	NoCollectedCards(ctx context.Context, chatID msginfo.ChatID) error
	ShowCollectedTotalInfo(
		ctx context.Context,
		chatID msginfo.ChatID,
		info []card.CollectedCardInfo,
	) error
	ShowCollectedReward(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		rew reward.Reward,
		count int,
		previousPage card.CollectedCardsPage,
		nextPage card.CollectedCardsPage,
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
