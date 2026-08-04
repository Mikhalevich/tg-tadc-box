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
	messageID msginfo.MessageID,
	rewardType reward.RewardType,
	page int,
) error {
	if page < 1 {
		return perror.InvalidParam("invalid page")
	}

	profile, _, err := v.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat id: %w", err)
	}

	maxPos := profile.Profile.Cards.CardsMaxPos(rewardType)

	if maxPos == 0 {
		if err := v.notifier.NoCollectedCards(ctx, chatID); err != nil {
			return fmt.Errorf("no collected cards notification: %w", err)
		}

		return nil
	}

	if page > maxPos {
		return perror.InvalidParam("invalid page")
	}

	cardByPos := profile.Profile.Cards.CardByPos(rewardType, page)

	collectedReward, err := v.rewardProvider.GetRewardByID(ctx, cardByPos.RewardID)
	if err != nil {
		return fmt.Errorf("get reward by id: %w", err)
	}

	if err := v.notifier.ShowCollectedReward(
		ctx,
		chatID,
		messageID,
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
