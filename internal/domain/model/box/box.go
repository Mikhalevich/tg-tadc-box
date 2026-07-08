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

type Type string

const (
	TypeNormal Type = "normal"
)

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
