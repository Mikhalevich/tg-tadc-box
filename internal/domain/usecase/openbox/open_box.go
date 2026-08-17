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
	GetBoxByID(ctx context.Context, id box.ID) (box.Box, error)
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

type BoxScheduler interface {
	ScheduleInProgressOrPending(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxType box.Type,
		meta box.Meta,
	) (box.Box, error)
	ActivatePending(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxType box.Type,
		now time.Time,
	) error
}

type Notifier interface {
	ShowReward(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		receivedReward reward.Reward,
		openingBox box.Box,
		withLikeButtons bool,
	) error
	ShowBonusBox(ctx context.Context, bonusBox box.Box) error
	ShowReadyToOpenBox(ctx context.Context, domBox box.Box) error
}

type TimeProvider interface {
	Now() time.Time
}

type OpenBox struct {
	bonusBoxAttempts map[box.Type]int
	repo             Repository
	transactor       Transactor
	playerProvider   PlayerProvider
	rewardGenerator  RewardGenerator
	boxScheduler     BoxScheduler
	notifier         Notifier
	timeProvider     TimeProvider
}

func New(
	bonusBoxAttempts map[box.Type]int,
	repo Repository,
	transactor Transactor,
	playerProvider PlayerProvider,
	rewardGenertor RewardGenerator,
	boxScheduler BoxScheduler,
	notifier Notifier,
	timeProvider TimeProvider,
) *OpenBox {
	return &OpenBox{
		bonusBoxAttempts: bonusBoxAttempts,
		repo:             repo,
		transactor:       transactor,
		playerProvider:   playerProvider,
		rewardGenerator:  rewardGenertor,
		boxScheduler:     boxScheduler,
		notifier:         notifier,
		timeProvider:     timeProvider,
	}
}
