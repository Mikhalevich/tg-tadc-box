package player

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
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

func (cc *CardsCollected) CardByPos(rewardType reward.RewardType, pos int) (Card, error) {
	cardRef, err := cc.cardByPosRef(rewardType, pos)
	if err != nil {
		return Card{}, fmt.Errorf("card by pos ref: %w", err)
	}

	return *cardRef, nil
}

func findCardIdxByRewardID(cards []Card, rewardID reward.ID) int {
	for idx, crd := range cards {
		if crd.RewardID == rewardID {
			return idx
		}
	}

	return -1
}

func (cc *CardsCollected) AbstractByPos(
	rewardType reward.RewardType,
	pos int,
	count int,
	abstractionCosts map[reward.RewardType]gloink.Amount,
) (gloink.Amount, error) {
	cardRef, err := cc.cardByPosRef(rewardType, pos)
	if err != nil {
		return 0, fmt.Errorf("card by pos ref: %w", err)
	}

	if cardRef.Count <= count {
		return 0, perror.InvalidParam("Not enaught cards to abstract")
	}

	cardRef.Count -= count

	return abstractionCosts[rewardType].Multiply(count), nil
}

// AbstractDuplicatesAll remove all duplicates from cards
// returns gloinks amount of abstracted cards.
func (cc *CardsCollected) AbstractDuplicatesAll(
	abstractionCosts map[reward.RewardType]gloink.Amount,
) gloink.Amount {
	return cc.processAbstractionDuplicatesAll(abstractionCosts, removeDuplicates)
}

// ViewCostOfAbstractionDuplicatesAll view possible gloinks amount
// if doing AbstractDuplicatesAll function.
func (cc *CardsCollected) ViewCostOfAbstractionDuplicatesAll(
	abstractionCosts map[reward.RewardType]gloink.Amount,
) gloink.Amount {
	return cc.processAbstractionDuplicatesAll(abstractionCosts, countDuplicates)
}

func (cc *CardsCollected) processAbstractionDuplicatesAll(
	abstractionCosts map[reward.RewardType]gloink.Amount,
	processorFn func(cards []Card) int,
) gloink.Amount {
	var gloinksAmount gloink.Amount
	for rewardType, cards := range cc.CardsByType {
		count := processorFn(cards)
		gloinksAmount += abstractionCosts[rewardType].Multiply(count)
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

func (cc *CardsCollected) cardByPosRef(rewardType reward.RewardType, pos int) (*Card, error) {
	cards := cc.CardsByType[rewardType]

	if len(cards) < pos {
		return nil, perror.InvalidParam("Position out of range")
	}

	return &cards[pos-1], nil
}
