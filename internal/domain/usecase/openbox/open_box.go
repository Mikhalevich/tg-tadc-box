package openbox

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

type Repository interface {
	GetBoxByID(id int) (box.Box, error)
	UpdateBox(b box.Box) error
}

type Notifier interface {
	ShowBoxInfo(b box.Box, availableAfter time.Duration)
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
