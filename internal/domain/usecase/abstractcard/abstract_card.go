package abstractcard

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type PlayerProvider interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
	UpdatePlayer(ctx context.Context, usr player.Player) error
}

type PageProvider interface {
	Page(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		rewardType reward.RewardType,
		page int,
	) error
}

type Notifier interface {
	ShowGloinksWalletAfterAbstraction(
		ctx context.Context,
		chatID msginfo.ChatID,
		wallet player.Wallet,
		abstractedAmount gloink.Amount,
	) error
}

type AbstractCard struct {
	abstractionCosts map[reward.RewardType]gloink.Amount
	transactor       Transactor
	playerProvider   PlayerProvider
	pageProvider     PageProvider
	notifier         Notifier
}

func New(
	abstractionCosts map[reward.RewardType]gloink.Amount,
	transactor Transactor,
	playerProvider PlayerProvider,
	pageProvider PageProvider,
	notifier Notifier,
) *AbstractCard {
	return &AbstractCard{
		abstractionCosts: abstractionCosts,
		transactor:       transactor,
		playerProvider:   playerProvider,
		pageProvider:     pageProvider,
		notifier:         notifier,
	}
}
