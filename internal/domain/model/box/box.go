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
	TypeCommon    Type = "common"
	TypeRare      Type = "rare"
	TypeEpic      Type = "epic"
	TypeLegendary Type = "legendary"
)

func (t Type) String() string {
	return string(t)
}

func (t Type) Compare(other Type) int {
	if t == other {
		return 0
	}

	for _, inOrder := range []Type{TypeCommon, TypeRare, TypeEpic, TypeLegendary} {
		if t == inOrder {
			return -1
		}

		if other == inOrder {
			return 1
		}
	}

	return 0
}

func (t Type) IsValid() bool {
	switch t {
	case TypeCommon, TypeRare, TypeEpic, TypeLegendary:
		return true
	}

	return false
}

func (t Type) Next() Type {
	switch t {
	case TypeCommon:
		return TypeRare

	case TypeRare:
		return TypeEpic

	case TypeEpic:
		return TypeLegendary

	case TypeLegendary:
	}

	return Type("")
}

func TypeFromString(str string) (Type, error) {
	switch str {
	case TypeCommon.String(), TypeRare.String(), TypeEpic.String(), TypeLegendary.String():
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
	Meta                Meta
}

func (b Box) IsStatus(statuses ...Status) bool {
	return slices.Contains(statuses, b.Status)
}

func (b Box) AvailableAfter(now time.Time) time.Duration {
	return b.AvailableAt.Sub(now)
}

type InProgressBox struct {
	Box            Box
	AvailableAfter time.Duration
}

func SortBoxByType(boxes []Box) {
	slices.SortFunc(boxes, func(a, b Box) int {
		return a.Type.Compare(b.Type)
	})
}
