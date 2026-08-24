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

type PlayerService interface {
	AbstractCard(
		ctx context.Context,
		chatID msginfo.ChatID,
		rewardType reward.RewardType,
		pos int,
		count int,
	) (gloink.Amount, player.Wallet, error)
	AbstractAllDuplicates(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (gloink.Amount, player.Wallet, error)
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
	transactor    Transactor
	playerService PlayerService
	pageProvider  PageProvider
	notifier      Notifier
}

func New(
	transactor Transactor,
	playerService PlayerService,
	pageProvider PageProvider,
	notifier Notifier,
) *AbstractCard {
	return &AbstractCard{
		transactor:    transactor,
		playerService: playerService,
		pageProvider:  pageProvider,
		notifier:      notifier,
	}
}
