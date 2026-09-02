package openbox

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type BoxService interface {
	GetBoxByID(ctx context.Context, id box.ID) (box.Box, error)
	OpenBox(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxID box.ID,
		receivedRewardID reward.ID,
		completedAt time.Time,
	) error
	ActivatePending(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxType box.Type,
		now time.Time,
	) error
	ScheduleInProgressOrPending(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxType box.Type,
		createdAt time.Time,
		meta box.Meta,
	) (box.Box, error)
}

type PlayerService interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, bool, error)
	UpdatePlayer(ctx context.Context, usr player.Player) error
}

type RewardService interface {
	Generate(ctx context.Context, boxType box.Type) (reward.Reward, error)
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
	ShowBonusBox(
		ctx context.Context,
		bonusBox box.Box,
	) error
}

type TimeProvider interface {
	Now() time.Time
}

type OpenBox struct {
	bonusBoxAttempts map[box.Type]int
	transactor       Transactor
	boxService       BoxService
	playerService    PlayerService
	rewardService    RewardService
	notifier         Notifier
	timeProvider     TimeProvider
}

func New(
	bonusBoxAttempts map[box.Type]int,
	transactor Transactor,
	boxService BoxService,
	playerService PlayerService,
	rewardService RewardService,
	notifier Notifier,
	timeProvider TimeProvider,
) *OpenBox {
	return &OpenBox{
		bonusBoxAttempts: bonusBoxAttempts,
		transactor:       transactor,
		boxService:       boxService,
		playerService:    playerService,
		rewardService:    rewardService,
		notifier:         notifier,
		timeProvider:     timeProvider,
	}
}
