package card

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type ID int

func (id ID) Int() int {
	return int(id)
}

func IDFromInt(id int) ID {
	return ID(id)
}

type Card struct {
	ID        ID
	ChatID    msginfo.ChatID
	RewardID  reward.ID
	Count     int
	UpdatedAt time.Time
}

type CollectedCardsPage struct {
	Page    int
	IsValid bool
}

type CollectedCardInfo struct {
	Type      reward.RewardType
	Collected int
	Total     int
}
