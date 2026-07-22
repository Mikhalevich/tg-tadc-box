package openbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (o *OpenBox) OpenByID(
	ctx context.Context,
	chatID msginfo.ChatID,
	id int,
) error {
	readyBox, err := o.repo.GetBoxByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get box by id: %w", err)
	}

	if err := isInProgress(readyBox.Status); err != nil {
		return fmt.Errorf("can be opened: %w", err)
	}

	now := o.timeProvider.Now()

	if readyBox.AvailableAt.After(now) {
		if err := o.notifier.ShowBoxInfo(ctx, readyBox, readyBox.AvailableAt.Sub(now)); err != nil {
			return fmt.Errorf("show box info: %w", err)
		}

		return nil
	}

	receivedReward, err := o.openBox(ctx, chatID, readyBox, now)
	if err != nil {
		return fmt.Errorf("open box: %w", err)
	}

	if err := o.notifier.ShowReward(ctx, chatID, receivedReward); err != nil {
		return fmt.Errorf("show reward: %w", err)
	}

	return nil
}

func (o *OpenBox) openBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	readyBox box.Box,
	completedAt time.Time,
) (reward.Reward, error) {
	readyBox.Status = box.StatusOpened
	readyBox.CompletedAt = completedAt

	if err := o.repo.UpdateBox(ctx, readyBox); err != nil {
		return reward.Reward{}, fmt.Errorf("update box: %w", err)
	}

	receivedReward, err := o.rewardGenerator.Generate(ctx)
	if err != nil {
		return reward.Reward{}, fmt.Errorf("generate reward: %w", err)
	}

	if err := o.repo.InsertReceivedReward(ctx, reward.ReceivedReward{
		ChatID:    chatID,
		RewardID:  receivedReward.ID,
		BoxID:     readyBox.ID,
		CreatedAt: completedAt,
	}); err != nil {
		return reward.Reward{}, fmt.Errorf("insert received reward: %w", err)
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
