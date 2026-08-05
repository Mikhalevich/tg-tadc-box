package shop

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type PlayerProvider interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
	UpdatePlayer(ctx context.Context, plr player.Player) error
}

type BoxScheduler interface {
	ScheduleBox(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxType box.Type,
	) error
}

type Notifier interface {
	ShowBoxCosts(
		ctx context.Context,
		chatID msginfo.ChatID,
		wallet player.Wallet,
		costs []gloink.BoxCost,
	) error
}

type Shop struct {
	boxCosts       []gloink.BoxCost
	transactor     Transactor
	playerProvider PlayerProvider
	boxScheduler   BoxScheduler
	notifier       Notifier
}

func New(
	boxCosts []gloink.BoxCost,
	transactor Transactor,
	playerProvider PlayerProvider,
	boxScheduler BoxScheduler,
	notifier Notifier,
) *Shop {
	return &Shop{
		boxCosts:       boxCosts,
		transactor:     transactor,
		playerProvider: playerProvider,
		boxScheduler:   boxScheduler,
		notifier:       notifier,
	}
}
