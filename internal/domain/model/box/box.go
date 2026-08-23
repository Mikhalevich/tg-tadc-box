package box

import (
	"fmt"
	"slices"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
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

func (t Type) Pretty() string {
	return cases.Title(language.English).String(string(t))
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

func (b Box) IsValid() bool {
	return b.ID.Int() != 0
}

func (b Box) IsStatus(statuses ...Status) bool {
	return slices.Contains(statuses, b.Status)
}

func (b Box) AvailableAfter(now time.Time) time.Duration {
	return b.AvailableAt.Sub(now)
}

func (b Box) IsInProgress() error {
	switch b.Status {
	case StatusPending:
		return perror.InvalidStatus("your box is in pending status")

	case StatusOpened:
		return perror.InvalidStatus("box already opened")

	case StatusCanceled:
		return perror.InvalidStatus("box already canceled")

	case StatusInProgress:
	}

	return nil
}

func (b Box) IsReadyToOpen(now time.Time) error {
	if err := b.IsInProgress(); err != nil {
		return fmt.Errorf("not is in_progress status: %w", err)
	}

	if b.AvailableAt.After(now) {
		return perror.InvalidParam("box is not ready yet")
	}

	return nil
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

func ToInProgressBoxMapByType(boxes []Box, now time.Time) map[Type]InProgressBox {
	if len(boxes) == 0 {
		return nil
	}

	boxesMap := make(map[Type]InProgressBox, len(boxes))

	for _, b := range boxes {
		boxesMap[b.Type] = InProgressBox{
			Box:            b,
			AvailableAfter: b.AvailableAfter(now),
		}
	}

	return boxesMap
}
