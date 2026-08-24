package viewcards

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type PlayerService interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
}

type RewardService interface {
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
	playerService    PlayerService
	rewardService    RewardService
	notifier         Notifier
}

func New(
	abstractionCosts map[reward.RewardType]gloink.Amount,
	playerService PlayerService,
	rewardService RewardService,
	notifier Notifier,
) *ViewCards {
	return &ViewCards{
		abstractionCosts: abstractionCosts,
		playerService:    playerService,
		rewardService:    rewardService,
		notifier:         notifier,
	}
}
