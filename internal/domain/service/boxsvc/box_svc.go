package boxsvc

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type Repository interface {
	GetBoxByID(ctx context.Context, id box.ID) (box.Box, error)
	CompleteBox(
		ctx context.Context,
		boxID box.ID,
		status box.Status,
		completedAt time.Time,
	) error
	UpdateBox(ctx context.Context, b box.Box) error
	InsertReceivedReward(
		ctx context.Context,
		rwd reward.ReceivedReward,
	) error
	ChangeFirstBoxStatusByType(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxType box.Type,
		newStatus box.Status,
		previousStatus box.Status,
		availableAt time.Time,
	) error
	GetBoxesByStatus(ctx context.Context, chatID msginfo.ChatID, statuses ...box.Status) ([]box.Box, error)
	InsertBox(ctx context.Context, b box.Box) (int, error)
}

type Service struct {
	boxWaitPeriod map[box.Type]time.Duration
	transactor    Transactor
	repo          Repository
}

func New(
	boxWaitPeriod map[box.Type]time.Duration,
	transactor Transactor,
	repo Repository,
) *Service {
	return &Service{
		boxWaitPeriod: boxWaitPeriod,
		transactor:    transactor,
		repo:          repo,
	}
}
