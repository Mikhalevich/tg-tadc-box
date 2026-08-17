package openbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
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
		profile, err := o.playerProvider.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get user by chat_id: %w", err)
		}

		readyBox, err := o.repo.GetBoxByID(ctx, boxID)
		if err != nil {
			return fmt.Errorf("get box by id: %w", err)
		}

		now := o.timeProvider.Now()

		if err := isBoxReadyToOpen(readyBox, now); err != nil {
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

		if err := o.boxScheduler.ActivatePending(
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
	readyBox.Status = box.StatusOpened
	readyBox.CompletedAt = completedAt

	if err := o.repo.UpdateBox(ctx, readyBox); err != nil {
		return reward.Reward{}, fmt.Errorf("update box: %w", err)
	}

	receivedReward, err := o.rewardGenerator.Generate(ctx, readyBox.Type)
	if err != nil {
		return reward.Reward{}, fmt.Errorf("generate reward: %w", err)
	}

	if err := o.repo.InsertReceivedReward(ctx, reward.ReceivedReward{
		ChatID:    profile.ChatID,
		RewardID:  receivedReward.ID,
		BoxID:     readyBox.ID,
		CreatedAt: completedAt,
	}); err != nil {
		return reward.Reward{}, fmt.Errorf("insert received reward: %w", err)
	}

	if err := o.markRewardInUserProfile(ctx, profile, receivedReward, readyBox.Type); err != nil {
		return reward.Reward{}, fmt.Errorf("mark opened box in user profile: %w", err)
	}

	return receivedReward, nil
}

func isInProgress(status box.Status) error {
	switch status {
	case box.StatusPending:
		return perror.InvalidStatus("your box is in pending status")

	case box.StatusOpened:
		return perror.InvalidStatus("box already opened")

	case box.StatusCanceled:
		return perror.InvalidStatus("box already canceled")

	case box.StatusInProgress:
	}

	return nil
}

func (o *OpenBox) markRewardInUserProfile(
	ctx context.Context,
	plr player.Player,
	rwd reward.Reward,
	boxType box.Type,
) error {
	plr.Profile.Cards.AddReward(rwd.Type, rwd.ID, rwd.CreatedAt)

	plr.Profile.OpenedBoxes.Add(boxType)

	if err := o.playerProvider.UpdatePlayer(ctx, plr); err != nil {
		return fmt.Errorf("update player: %w", err)
	}

	if err := o.assignBonusBoxIfAvailable(
		ctx,
		plr,
		boxType,
	); err != nil {
		return fmt.Errorf("assign free box if available: %w", err)
	}

	return nil
}

func (o *OpenBox) assignBonusBoxIfAvailable(
	ctx context.Context,
	plr player.Player,
	openedBoxType box.Type,
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

	bonusBox, err := o.boxScheduler.ScheduleInProgressOrPending(
		ctx,
		plr.ChatID,
		bonusBoxType,
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
