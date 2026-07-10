package openbox

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type Repository interface {
	GetBoxesByStatus(ctx context.Context, chatID msginfo.ChatID, statuses ...box.Status) ([]box.Box, error)
	InsertBox(ctx context.Context, b box.Box) (int, error)
	GetBoxByID(ctx context.Context, id int) (box.Box, error)
	UpdateBox(ctx context.Context, b box.Box) error
}

type Notifier interface {
	ShowBoxInfo(ctx context.Context, b box.Box, availableAfter time.Duration) error
	ShowReward(ctx context.Context, chatID msginfo.ChatID) error
}

type TimeProvider interface {
	Now() time.Time
}

type OpenBox struct {
	repo         Repository
	notifier     Notifier
	timeProvider TimeProvider
}

func New(
	repo Repository,
	notifier Notifier,
	timeProvider TimeProvider,
) *OpenBox {
	return &OpenBox{
		repo:         repo,
		notifier:     notifier,
		timeProvider: timeProvider,
	}
}
