package schedulebox

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type Repository interface {
	GetActiveBoxes(chatID msginfo.ChatID) ([]box.Box, error)
	InsertBox(b box.Box) (int, error)
}

type Notifier interface {
	ShowBoxInfo(b box.Box, availableAfter time.Duration)
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
