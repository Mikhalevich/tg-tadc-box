package openbox

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
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
	InsertCard(
		ctx context.Context,
		receivedCard card.Card,
	) (int, error)
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type RewardGenerator interface {
	Generate(ctx context.Context) (reward.Reward, error)
}

type Notifier interface {
	ShowBoxInfo(ctx context.Context, b box.Box, availableAfter time.Duration) error
	ShowReward(ctx context.Context, chatID msginfo.ChatID, receivedReward reward.Reward) error
}

type TimeProvider interface {
	Now() time.Time
}

type OpenBox struct {
	repo            Repository
	transacor       Transactor
	rewardGenerator RewardGenerator
	notifier        Notifier
	timeProvider    TimeProvider
}

func New(
	repo Repository,
	transactor Transactor,
	rewardGenertor RewardGenerator,
	notifier Notifier,
	timeProvider TimeProvider,
) *OpenBox {
	return &OpenBox{
		repo:            repo,
		transacor:       transactor,
		rewardGenerator: rewardGenertor,
		notifier:        notifier,
		timeProvider:    timeProvider,
	}
}
