package viewcards

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (v *ViewCards) Page(
	ctx context.Context,
	chatID msginfo.ChatID,
	rewardType reward.RewardType,
	page int,
) error {
	if page < 1 {
		return perror.InvalidParam("invalid page")
	}

	maxPos, err := v.repo.GetCollectedCardMaxPos(ctx, chatID, rewardType)
	if err != nil {
		return fmt.Errorf("get collected card max pos: %w", err)
	}

	if maxPos == 0 {
		if err := v.notifier.NoCollectedCards(ctx, chatID); err != nil {
			return fmt.Errorf("no collected cards notification: %w", err)
		}

		return nil
	}

	if page > maxPos {
		return perror.InvalidParam("invalid page")
	}

	cardByPos, err := v.repo.GetCollectedCardByPos(ctx, chatID, rewardType, page)
	if err != nil {
		return fmt.Errorf("get collected card by pos %d: %w", page, err)
	}

	collectedReward, err := v.rewardProvider.GetRewardByID(ctx, cardByPos.RewardID)
	if err != nil {
		return fmt.Errorf("get reward by id: %w", err)
	}

	if err := v.notifier.ShowCollectedReward(
		ctx,
		chatID,
		collectedReward,
		cardByPos.Count,
		pageInfo(page-1, maxPos),
		pageInfo(page+1, maxPos),
	); err != nil {
		return fmt.Errorf("show collected card: %w", err)
	}

	return nil
}

func pageInfo(page, maxPage int) card.CollectedCardsPage {
	if page < 1 || page > maxPage {
		return card.CollectedCardsPage{}
	}

	return card.CollectedCardsPage{
		Page:    page,
		IsValid: true,
	}
}
