package model

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Card struct {
	ID        int       `db:"id"`
	ChatID    int64     `db:"chat_id"`
	RewardID  int       `db:"reward_id"`
	Count     int       `db:"count"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c Card) ToDomCard() card.Card {
	return card.Card{
		ID:        card.IDFromInt(c.ID),
		ChatID:    msginfo.ChatIDFromInt64(c.ChatID),
		RewardID:  reward.IDFromInt(c.RewardID),
		Count:     c.Count,
		UpdatedAt: c.UpdatedAt,
	}
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

type CollectedCardTotalInfo struct {
	Type  string `db:"type"`
	Total int    `db:"total"`
}
