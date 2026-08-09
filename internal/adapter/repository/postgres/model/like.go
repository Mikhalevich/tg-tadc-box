package model

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/like"
)

type Like struct {
	ID        int       `db:"id"`
	BoxID     int       `db:"box_id"`
	ChatID    int64     `db:"chat_id"`
	RewardID  int       `db:"reward_id"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
}

func ToDBLike(domLike like.Like) Like {
	return Like{
		ID:        domLike.ID,
		BoxID:     domLike.BoxID.Int(),
		ChatID:    domLike.ChatID.Int64(),
		RewardID:  domLike.RewardID.Int(),
		Type:      domLike.Type.String(),
		CreatedAt: domLike.CreatedAt,
	}
}
