package openbox

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

type Repository interface {
	GetBoxByID(ctx context.Context, id int) (box.Box, error)
	UpdateBox(ctx context.Context, b box.Box) error
}

type Notifier interface {
	ShowBoxInfo(ctx context.Context, b box.Box, availableAfter time.Duration) error
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
