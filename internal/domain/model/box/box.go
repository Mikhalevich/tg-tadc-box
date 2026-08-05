package box

import (
	"fmt"
	"slices"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type ID int

func (id ID) Int() int {
	return int(id)
}

func IDFromInt(id int) ID {
	return ID(id)
}

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
	TypeNormal    Type = "normal"
	TypeRare      Type = "rare"
	TypeEpic      Type = "epic"
	TypeLegendary Type = "legendary"
)

func (t Type) String() string {
	return string(t)
}

func TypeFromString(str string) (Type, error) {
	switch str {
	case TypeNormal.String(), TypeRare.String(), TypeEpic.String(), TypeLegendary.String():
		return Type(str), nil
	}

	return "", fmt.Errorf("invalid type string %q", str)
}

type Box struct {
	ID                  ID
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
