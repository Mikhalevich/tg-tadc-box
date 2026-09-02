package openbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (o *OpenBox) OpenByID(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	boxID box.ID,
) error {
	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		profile, _, err := o.playerService.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get user by chat_id: %w", err)
		}

		readyBox, err := o.boxService.GetBoxByID(ctx, boxID)
		if err != nil {
			return fmt.Errorf("get box by id: %w", err)
		}

		now := o.timeProvider.Now()

		if err := readyBox.IsReadyToOpen(now); err != nil {
			return fmt.Errorf("check box is ready for open: %w", err)
		}

		receivedReward, err := o.openBox(ctx, profile, readyBox, now)
		if err != nil {
			return fmt.Errorf("open box: %w", err)
		}

		if err := o.notifier.ShowReward(
			ctx,
			chatID,
			messageID,
			receivedReward,
			readyBox,
			true,
		); err != nil {
			return fmt.Errorf("show reward: %w", err)
		}

		if err := o.boxService.ActivatePending(
			ctx,
			chatID,
			readyBox.Type,
			now,
		); err != nil {
			return fmt.Errorf("activate pending: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (o *OpenBox) openBox(
	ctx context.Context,
	profile player.Player,
	readyBox box.Box,
	completedAt time.Time,
) (reward.Reward, error) {
	receivedReward, err := o.rewardService.Generate(ctx, readyBox.Type)
	if err != nil {
		return reward.Reward{}, fmt.Errorf("generate reward: %w", err)
	}

	if err := o.boxService.OpenBox(
		ctx,
		profile.ChatID,
		readyBox.ID,
		receivedReward.ID,
		completedAt,
	); err != nil {
		return reward.Reward{}, fmt.Errorf("open box: %w", err)
	}

	if err := o.markRewardInUserProfile(ctx, profile, receivedReward, readyBox.Type, completedAt); err != nil {
		return reward.Reward{}, fmt.Errorf("mark opened box in user profile: %w", err)
	}

	return receivedReward, nil
}

func (o *OpenBox) markRewardInUserProfile(
	ctx context.Context,
	plr player.Player,
	rwd reward.Reward,
	boxType box.Type,
	createdAt time.Time,
) error {
	plr.Profile.Cards.AddReward(rwd.Type, rwd.ID, rwd.CreatedAt)

	plr.Profile.OpenedBoxes.Add(boxType)

	if err := o.playerService.UpdatePlayer(ctx, plr); err != nil {
		return fmt.Errorf("update player: %w", err)
	}

	if err := o.assignBonusBoxIfAvailable(
		ctx,
		plr,
		boxType,
		createdAt,
	); err != nil {
		return fmt.Errorf("assign free box if available: %w", err)
	}

	return nil
}

func (o *OpenBox) assignBonusBoxIfAvailable(
	ctx context.Context,
	plr player.Player,
	openedBoxType box.Type,
	createdAt time.Time,
) error {
	count, ok := o.bonusBoxAttempts[openedBoxType]
	if !ok {
		return nil
	}

	bonusBoxType, isAvailable := plr.Profile.OpenedBoxes.IsBonusBoxAvailable(
		openedBoxType,
		count,
	)
	if !isAvailable {
		return nil
	}

	bonusBox, err := o.boxService.ScheduleInProgressOrPending(
		ctx,
		plr.ChatID,
		bonusBoxType,
		createdAt,
		box.Meta{
			BonusBox: box.BonusBox{
				IsValid:  true,
				Type:     openedBoxType,
				Attempts: count,
			},
		},
	)

	if err != nil {
		return fmt.Errorf("schedule bonus box: %w", err)
	}

	if err := o.notifier.ShowBonusBox(
		ctx,
		bonusBox,
	); err != nil {
		return fmt.Errorf("show bonus box: %w", err)
	}

	return nil
}
