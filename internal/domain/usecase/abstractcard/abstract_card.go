package abstractcard

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

type PlayerProvider interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
	UpdatePlayer(ctx context.Context, usr player.Player) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
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
	playerProvider PlayerProvider
	transactor     Transactor
	notifier       Notifier
}

func New(
	playerProvider PlayerProvider,
	transactor Transactor,
	notifier Notifier,
) *AbstractCard {
	return &AbstractCard{
		playerProvider: playerProvider,
		transactor:     transactor,
		notifier:       notifier,
	}
}
