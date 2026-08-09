package likereward

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/like"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type Repository interface {
	InsertLike(
		ctx context.Context,
		rwdLike like.Like,
	) error

	GetBoxByID(
		ctx context.Context,
		boxID box.ID,
	) (box.Box, error)

	GetRewardByID(
		ctx context.Context,
		rewardID reward.ID,
	) (reward.Reward, error)
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
}

type TimeProvider interface {
	Now() time.Time
}

type LikeReward struct {
	transactor   Transactor
	repo         Repository
	notifier     Notifier
	timeProvider TimeProvider
}

func New(
	transactor Transactor,
	repo Repository,
	notifier Notifier,
	timeProvider TimeProvider,
) *LikeReward {
	return &LikeReward{
		transactor:   transactor,
		repo:         repo,
		notifier:     notifier,
		timeProvider: timeProvider,
	}
}
