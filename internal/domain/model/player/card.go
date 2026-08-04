package player

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

var (
	//nolint:gochecknoglobals,mnd
	abstractionGloinksCosts = map[reward.RewardType]int{
		reward.RewardTypeCommon:    1,
		reward.RewardTypeRare:      5,
		reward.RewardTypeEpic:      25,
		reward.RewardTypeLegendary: 100,
	}
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

// AbstractDuplicatesAll remove all duplicates from cards
// returns gloinks amount of abstracted cards.
func (cc *CardsCollected) AbstractDuplicatesAll() int {
	return cc.processAbstractionDuplicatesAll(removeDuplicates)
}

// ViewCostOfAbstractionDuplicatesAll view possible gloinks amount
// if doing AbstractDuplicatesAll function.
func (cc *CardsCollected) ViewCostOfAbstractionDuplicatesAll() int {
	return cc.processAbstractionDuplicatesAll(countDuplicates)
}

func (cc *CardsCollected) processAbstractionDuplicatesAll(
	processorFn func(cards []Card) int,
) int {
	gloinksAmount := 0
	for rewardType, cards := range cc.CardsByType {
		count := processorFn(cards)
		gloinksAmount += abstractionGloinksCosts[rewardType] * count
	}

	return gloinksAmount
}

// removeDuplicates removes duplicates
// returns count of removing cards.
func removeDuplicates(cards []Card) int {
	removedCount := 0
	for idx, crd := range cards {
		if crd.Count > 1 {
			cards[idx].Count = 1
			removedCount += crd.Count - 1
		}
	}

	return removedCount
}

// countDuplicates counts duplicates.
func countDuplicates(cards []Card) int {
	duplicatesCount := 0
	for _, crd := range cards {
		if crd.Count > 1 {
			duplicatesCount += crd.Count - 1
		}
	}

	return duplicatesCount
}
