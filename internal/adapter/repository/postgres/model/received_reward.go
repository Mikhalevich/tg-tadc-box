package model

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type ReceivedReward struct {
	ID        int       `db:"id"`
	ChatID    int64     `db:"chat_id"`
	RewardID  int       `db:"reward_id"`
	BoxID     int       `db:"box_id"`
	CreatedAt time.Time `db:"created_at"`
}

func ToDBReceivedReward(domReward reward.ReceivedReward) ReceivedReward {
	return ReceivedReward{
		ID:        domReward.ID.Int(),
		ChatID:    domReward.ChatID.Int64(),
		RewardID:  domReward.RewardID.Int(),
		BoxID:     domReward.BoxID.Int(),
		CreatedAt: domReward.CreatedAt,
	}
}
