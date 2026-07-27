package player

import (
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

type Player struct {
	ID               ID
	ChatID           msginfo.ChatID
	CreatedAt        time.Time
	Profile          Profile
	ProfileVersion   int
	ProfileUpdatedAt time.Time
}
