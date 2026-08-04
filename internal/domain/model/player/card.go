package player

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Card struct {
	RewardID  reward.ID
	Count     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CardsCollected struct {
	CardsByType map[reward.RewardType][]Card
}

func (cc *CardsCollected) AddReward(
	rewardType reward.RewardType,
	rewardID reward.ID,
	createdAt time.Time,
) {
	if cc.CardsByType == nil {
		cc.CardsByType = make(map[reward.RewardType][]Card)
	}

	collectedCards, ok := cc.CardsByType[rewardType]
	if !ok {
		cc.CardsByType[rewardType] = []Card{
			{
				RewardID:  rewardID,
				Count:     1,
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			},
		}

		return
	}

	if collectedCardIdx := findCardIdxByRewardID(collectedCards, rewardID); collectedCardIdx != -1 {
		collectedCards[collectedCardIdx].Count++
		collectedCards[collectedCardIdx].UpdatedAt = createdAt

		return
	}

	cc.CardsByType[rewardType] = append(cc.CardsByType[rewardType], Card{
		RewardID:  rewardID,
		Count:     1,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	})
}

func (cc *CardsCollected) CardsCount() map[reward.RewardType]int {
	return map[reward.RewardType]int{
		reward.RewardTypeCommon:    len(cc.CardsByType[reward.RewardTypeCommon]),
		reward.RewardTypeRare:      len(cc.CardsByType[reward.RewardTypeRare]),
		reward.RewardTypeEpic:      len(cc.CardsByType[reward.RewardTypeEpic]),
		reward.RewardTypeLegendary: len(cc.CardsByType[reward.RewardTypeLegendary]),
	}
}

func (cc *CardsCollected) CardsMaxPos(rewardType reward.RewardType) int {
	return len(cc.CardsByType[rewardType])
}

func (cc *CardsCollected) CardByPos(rewardType reward.RewardType, pos int) Card {
	return cc.CardsByType[rewardType][pos-1]
}

func findCardIdxByRewardID(cards []Card, rewardID reward.ID) int {
	for idx, crd := range cards {
		if crd.RewardID == rewardID {
			return idx
		}
	}

	return -1
}
