package box

import (
	"slices"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusOpened     Status = "opened"
	StatusCanceled   Status = "canceled"
)

func (s Status) String() string {
	return string(s)
}

type Type string

const (
	TypeNormal Type = "normal"
)

func (t Type) String() string {
	return string(t)
}

type Box struct {
	ID          int
	ChatID      msginfo.ChatID
	Status      Status
	Type        Type
	CreatedAt   time.Time
	AvailableAt time.Time
	CompletedAt time.Time
}

func (b Box) IsStatus(statuses ...Status) bool {
	return slices.Contains(statuses, b.Status)
}

func (b Box) AvailableAfter(now time.Time) time.Duration {
	return b.AvailableAt.Sub(now)
}
