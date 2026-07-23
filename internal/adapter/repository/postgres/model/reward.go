package model

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Reward struct {
	ID        int       `db:"id"`
	Type      string    `db:"type"`
	URI       string    `db:"uri"`
	CreatedAt time.Time `db:"created_at"`
}

func (r Reward) ToDom() reward.Reward {
	return reward.Reward{
		ID:        reward.IDFromInt(r.ID),
		Type:      reward.RewardType(r.Type),
		URI:       r.URI,
		CreatedAt: r.CreatedAt,
	}
}

func ToDomRewards(dbRewards []Reward) []reward.Reward {
	domRewards := make([]reward.Reward, 0, len(dbRewards))

	for _, dbReward := range dbRewards {
		domRewards = append(domRewards, dbReward.ToDom())
	}

	return domRewards
}

type RewardTotalInfo struct {
	Type  string `db:"type"`
	Total int    `db:"total"`
}
