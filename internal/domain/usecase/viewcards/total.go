package viewcards

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (v *ViewCards) Total(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	profile, err := v.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat id: %w", err)
	}

	totalRewardCount, err := v.rewardProvider.GetRewardCountByType(ctx)
	if err != nil {
		return fmt.Errorf("get reward count by type: %w", err)
	}

	if err := v.notifier.ShowCollectedTotalInfo(
		ctx,
		chatID,
		makeCollectedCardsInfo(profile.Profile.Cards.CardsCount(), totalRewardCount),
		profile.Profile.Cards.ViewCostOfAbstractionDuplicatesAll(),
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
