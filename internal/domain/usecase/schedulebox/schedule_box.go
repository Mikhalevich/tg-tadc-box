package schedulebox

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type Repository interface {
	GetBoxesByStatus(ctx context.Context, chatID msginfo.ChatID, statuses ...box.Status) ([]box.Box, error)
	InsertBox(ctx context.Context, b box.Box) (int, error)
	ChangeFirstBoxStatusByType(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxType box.Type,
		newStatus box.Status,
		previousStatus box.Status,
		availableAt time.Time,
	) error
}

type TimeProvider interface {
	Now() time.Time
}

type ScheduleBox struct {
	boxWaitPeriod map[box.Type]time.Duration
	transactor    Transactor
	repo          Repository
	timeProvider  TimeProvider
}

func New(
	boxWaitPeriod map[box.Type]time.Duration,
	transactor Transactor,
	repo Repository,
	timeProvider TimeProvider,
) *ScheduleBox {
	return &ScheduleBox{
		boxWaitPeriod: boxWaitPeriod,
		transactor:    transactor,
		repo:          repo,
		timeProvider:  timeProvider,
	}
}
