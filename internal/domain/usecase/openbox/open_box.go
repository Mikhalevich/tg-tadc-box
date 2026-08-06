package openbox

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Repository interface {
	GetBoxesByStatus(ctx context.Context, chatID msginfo.ChatID, statuses ...box.Status) ([]box.Box, error)
	InsertBox(ctx context.Context, b box.Box) (int, error)
	GetBoxByID(ctx context.Context, id int) (box.Box, error)
	UpdateBox(ctx context.Context, b box.Box) error
	InsertReceivedReward(
		ctx context.Context,
		rwd reward.ReceivedReward,
	) error
}

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

type RewardGenerator interface {
	Generate(ctx context.Context, boxType box.Type) (reward.Reward, error)
}

type Notifier interface {
	ShowReward(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		receivedReward reward.Reward,
		boxType box.Type,
	) error
	ShowInProgressBoxes(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxes []box.InProgressBox,
	) error
	ShowBonusBox(ctx context.Context, bonusBox box.Box) error
}

type TimeProvider interface {
	Now() time.Time
}

type OpenBox struct {
	boxWaitPeriod   map[box.Type]time.Duration
	repo            Repository
	transactor      Transactor
	playerProvider  PlayerProvider
	rewardGenerator RewardGenerator
	notifier        Notifier
	timeProvider    TimeProvider
}

func New(
	boxWaitPeriod map[box.Type]time.Duration,
	repo Repository,
	transactor Transactor,
	playerProvider PlayerProvider,
	rewardGenertor RewardGenerator,
	notifier Notifier,
	timeProvider TimeProvider,
) *OpenBox {
	return &OpenBox{
		boxWaitPeriod:   boxWaitPeriod,
		repo:            repo,
		transactor:      transactor,
		playerProvider:  playerProvider,
		rewardGenerator: rewardGenertor,
		notifier:        notifier,
		timeProvider:    timeProvider,
	}
}
