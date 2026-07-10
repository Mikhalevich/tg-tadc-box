package readybox

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

type Repository interface {
	GetReadyToOpenBoxes(ctx context.Context, now time.Time, limit int) ([]box.Box, error)
	SetBoxReadyNotificationAt(ctx context.Context, ids []int, notificationAt time.Time) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type Notifier interface {
	ShowBoxInfo(ctx context.Context, b box.Box, availableAfter time.Duration) error
}

type TimeProvider interface {
	Now() time.Time
}

type ReadyBox struct {
	repo         Repository
	transactor   Transactor
	notifier     Notifier
	timeProvider TimeProvider
}

func New(
	repo Repository,
	transactor Transactor,
	notifier Notifier,
	timeProvider TimeProvider,
) *ReadyBox {
	return &ReadyBox{
		repo:         repo,
		transactor:   transactor,
		notifier:     notifier,
		timeProvider: timeProvider,
	}
}
