package schedulebox

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type Repository interface {
	GetBoxesByStatus(ctx context.Context, chatID msginfo.ChatID, statuses ...box.Status) ([]box.Box, error)
	InsertBox(ctx context.Context, b box.Box) (int, error)
}

type Notifier interface {
	ShowBoxInfo(ctx context.Context, b box.Box, availableAfter time.Duration) error
}

type TimeProvider interface {
	Now() time.Time
}

type ScheduleBox struct {
	repo         Repository
	notifier     Notifier
	timeProvider TimeProvider
}

func New(
	repo Repository,
	notifier Notifier,
	timeProvider TimeProvider,
) *ScheduleBox {
	return &ScheduleBox{
		repo:         repo,
		timeProvider: timeProvider,
	}
}
