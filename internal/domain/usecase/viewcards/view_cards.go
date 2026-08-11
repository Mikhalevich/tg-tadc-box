package viewcards

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type PlayerProvider interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
}

type RewardProvider interface {
	GetRewardByID(ctx context.Context, id reward.ID) (reward.Reward, error)
	GetRewardCountByType(ctx context.Context) (map[reward.RewardType]int, error)
}

type Notifier interface {
	NoCollectedCards(ctx context.Context, chatID msginfo.ChatID) error
	ShowCardPageTotal(
		ctx context.Context,
		chatID msginfo.ChatID,
		info []card.CardPageTotal,
		abstractDuplicatesAmount gloink.Amount,
	) error
	ShowCardPage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		rew reward.Reward,
		count int,
		page int,
		maxPage int,
		firstPage card.CardPage,
		previousPage card.CardPage,
		nextPage card.CardPage,
		lastPage card.CardPage,
	) error
}

type ViewCards struct {
	abstractionCosts map[reward.RewardType]gloink.Amount
	playerProvider   PlayerProvider
	rewardProvider   RewardProvider
	notifier         Notifier
}

func New(
	abstractionCosts map[reward.RewardType]gloink.Amount,
	playerProvider PlayerProvider,
	rewardProvider RewardProvider,
	notifier Notifier,
) *ViewCards {
	return &ViewCards{
		abstractionCosts: abstractionCosts,
		playerProvider:   playerProvider,
		rewardProvider:   rewardProvider,
		notifier:         notifier,
	}
}
