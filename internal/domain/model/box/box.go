package box

import (
	"fmt"
	"slices"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
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
	ID                  int
	ChatID              msginfo.ChatID
	Status              Status
	Type                Type
	CreatedAt           time.Time
	AvailableAt         time.Time
	ReadyNotificationAt time.Time
	CompletedAt         time.Time
}

func (b Box) IsStatus(statuses ...Status) bool {
	return slices.Contains(statuses, b.Status)
}

func (b Box) AvailableAfter(now time.Time) time.Duration {
	return b.AvailableAt.Sub(now)
}

type OpenBoxButtonPayload struct {
	ID int
}

func OpenBoxButton(boxID int) (button.Button, error) {
	btn, err := button.CreateButton(
		"Open",
		button.OperationOpenBox,
		true,
		OpenBoxButtonPayload{
			ID: boxID,
		},
	)

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}
