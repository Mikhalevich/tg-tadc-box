package viewcards

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (v *ViewCards) Total(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	collectedCardCount, err := v.repo.GetCollectedCardCountByType(ctx, chatID)
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			if err := v.notifier.NoCollectedCards(ctx, chatID); err != nil {
				return fmt.Errorf("no collected cards notification: %w", err)
			}

			return nil
		}

		return fmt.Errorf("get collected card count by type: %w", err)
	}

	totalRewardCount, err := v.rewardProvider.GetRewardCountByType(ctx)
	if err != nil {
		return fmt.Errorf("get reward count by type: %w", err)
	}

	if err := v.notifier.ShowCollectedTotalInfo(
		ctx,
		chatID,
		makeCollectedCardsInfo(collectedCardCount, totalRewardCount),
	); err != nil {
		return fmt.Errorf("show collected total info: %w", err)
	}

	return nil
}

func makeCollectedCardsInfo(collected, total map[reward.RewardType]int) []card.CollectedCardInfo {
	info := []card.CollectedCardInfo{
		{
			Type: reward.RewardTypeCommon,
		},
		{
			Type: reward.RewardTypeRare,
		},
		{
			Type: reward.RewardTypeEpic,
		},
		{
			Type: reward.RewardTypeLegendary,
		},
	}

	for idx, value := range info {
		info[idx].Collected = collected[value.Type]
		info[idx].Total = total[value.Type]
	}

	return info
}
