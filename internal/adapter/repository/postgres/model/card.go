package model

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
)

type Card struct {
	ID        int       `db:"id"`
	ChatID    int64     `db:"chat_id"`
	RewardID  int       `db:"reward_id"`
	Count     int       `db:"count"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ToDBCard(domCard card.Card) Card {
	return Card{
		ID:        domCard.ID.Int(),
		ChatID:    domCard.ChatID.Int64(),
		RewardID:  domCard.RewardID.Int(),
		Count:     domCard.Count,
		UpdatedAt: domCard.UpdatedAt,
	}
}
