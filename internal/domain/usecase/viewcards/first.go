package viewcards

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (v *ViewCards) First(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	nextCard, err := v.repo.GetCollectedCardAfterID(ctx, 0)
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			if err := v.notifier.NoCollectedCards(ctx, chatID); err != nil {
				return fmt.Errorf("no collected cards notification: %w", err)
			}

			return nil
		}

		return fmt.Errorf("get collected card after id: %w", err)
	}

	collectedReward, err := v.rewardProvider.GetRewardByID(ctx, nextCard.RewardID)
	if err != nil {
		return fmt.Errorf("get reward by id: %w", err)
	}

	if err := v.notifier.ShowCollectedReward(
		ctx,
		chatID,
		collectedReward,
		nextCard.Count,
	); err != nil {
		return fmt.Errorf("show collected card: %w", err)
	}

	return nil
}
