package likereward

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/like"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (lr *LikeReward) Like(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	boxID box.ID,
	rewardID reward.ID,
	likeType like.Type,
) error {
	if err := lr.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := lr.repo.InsertLike(
			ctx,
			like.Like{
				BoxID:     boxID,
				ChatID:    chatID,
				RewardID:  rewardID,
				Type:      likeType,
				CreatedAt: lr.timeProvider.Now(),
			}); err != nil {
			return fmt.Errorf("insert like: %w", err)
		}

		likeBox, err := lr.repo.GetBoxByID(ctx, boxID)
		if err != nil {
			return fmt.Errorf("get box by id: %w", err)
		}

		if likeBox.ChatID != chatID {
			return perror.InvalidParam("invalid user for box")
		}

		rwd, err := lr.repo.GetRewardByID(ctx, rewardID)
		if err != nil {
			return fmt.Errorf("get reward by id: %w", err)
		}

		if err := lr.notifier.ShowReward(
			ctx,
			chatID,
			messageID,
			rwd,
			likeBox,
			false,
		); err != nil {
			return fmt.Errorf("show reward: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
