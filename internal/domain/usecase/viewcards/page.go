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

	profile, err := v.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat id: %w", err)
	}

	maxPage := profile.Profile.Cards.CardsMaxPos(rewardType)

	if maxPage == 0 {
		if err := v.notifier.NoCollectedCards(ctx, chatID); err != nil {
			return fmt.Errorf("no collected cards notification: %w", err)
		}

		return nil
	}

	if page > maxPage {
		return perror.InvalidParam("invalid page")
	}

	cardByPos := profile.Profile.Cards.CardByPos(rewardType, page)

	collectedReward, err := v.rewardProvider.GetRewardByID(ctx, cardByPos.RewardID)
	if err != nil {
		return fmt.Errorf("get reward by id: %w", err)
	}

	if err := v.notifier.ShowCardPage(
		ctx,
		chatID,
		messageID,
		collectedReward,
		cardByPos.Count,
		page,
		maxPage,
		makeFirstPage(page),
		makeCardPage(page-1, maxPage),
		makeCardPage(page+1, maxPage),
		makeLastPage(page, maxPage),
	); err != nil {
		return fmt.Errorf("show collected card: %w", err)
	}

	return nil
}

func makeFirstPage(page int) card.CardPage {
	if page <= 1 {
		return card.CardPage{}
	}

	return card.CardPage{
		Page:    1,
		IsValid: true,
	}
}

func makeLastPage(page, maxPage int) card.CardPage {
	if page >= maxPage {
		return card.CardPage{}
	}

	return card.CardPage{
		Page:    maxPage,
		IsValid: true,
	}
}

func makeCardPage(page, maxPage int) card.CardPage {
	if page < 1 || page > maxPage {
		return card.CardPage{}
	}

	return card.CardPage{
		Page:    page,
		IsValid: true,
	}
}
