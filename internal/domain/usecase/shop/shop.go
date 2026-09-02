package shop

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type PlayerService interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, bool, error)
	UpdatePlayer(ctx context.Context, plr player.Player) error
}

type BoxService interface {
	ScheduleInProgress(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxType box.Type,
		createdAt time.Time,
		isImmediate bool,
	) (box.Box, error)
	GetBoxesByStatus(
		ctx context.Context,
		chatID msginfo.ChatID,
		statuses ...box.Status,
	) ([]box.Box, error)
}

type Notifier interface {
	ShowBoxCosts(
		ctx context.Context,
		chatID msginfo.ChatID,
		wallet player.Wallet,
		costs []gloink.BoxCost,
		inProgressBoxes map[box.Type]box.InProgressBox,
	) error
	ShowReadyToOpenBox(ctx context.Context, domBox box.Box) error
}

type TimeProvider interface {
	Now() time.Time
}

type Shop struct {
	boxCosts      []gloink.BoxCost
	transactor    Transactor
	playerService PlayerService
	boxService    BoxService
	notifier      Notifier
	timeProvider  TimeProvider
}

func New(
	boxCosts []gloink.BoxCost,
	transactor Transactor,
	playerService PlayerService,
	boxService BoxService,
	notifier Notifier,
	timeProvider TimeProvider,
) *Shop {
	return &Shop{
		boxCosts:      boxCosts,
		transactor:    transactor,
		playerService: playerService,
		boxService:    boxService,
		notifier:      notifier,
		timeProvider:  timeProvider,
	}
}
