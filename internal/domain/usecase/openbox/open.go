package openbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (o *OpenBox) Open(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		profile, err := o.playerProvider.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get user profile: %w", err)
		}

		now := o.timeProvider.Now()

		isAvailable, err := o.isNewBoxScheduleAvailable(ctx, chatID, now)
		if err != nil {
			return fmt.Errorf("is box schedule available: %w", err)
		}

		if !isAvailable {
			return nil
		}

		if err := o.scheduleBox(ctx, chatID, now, isNoOpenBoxes(profile)); err != nil {
			return fmt.Errorf("schedule box: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func isNoOpenBoxes(profile player.Player) bool {
	return profile.Profile.OpenedBoxes.Count() == 0
}

func (o *OpenBox) isNewBoxScheduleAvailable(
	ctx context.Context,
	chatID msginfo.ChatID,
	now time.Time,
) (bool, error) {
	boxes, err := o.repo.GetBoxesByStatus(ctx, chatID, box.StatusInProgress)
	if err != nil {
		return false, fmt.Errorf("get active boxes: %w", err)
	}

	if len(boxes) > 0 {
		if err := o.notifier.ShowBoxInfo(ctx, boxes[0], boxes[0].AvailableAfter(now)); err != nil {
			return false, fmt.Errorf("show existing box info: %w", err)
		}

		return false, nil
	}

	return true, nil
}

func (o *OpenBox) scheduleBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	createdAt time.Time,
	isImmediate bool,
) error {
	newBox := createNormalBox(
		chatID,
		createdAt,
		calculateWaitPeriod(o.commonWaitPeriod, isImmediate),
	)

	boxID, err := o.repo.InsertBox(ctx, newBox)
	if err != nil {
		return fmt.Errorf("insert box: %w", err)
	}

	newBox.ID = box.IDFromInt(boxID)

	if err := o.notifier.ShowBoxInfo(ctx, newBox, newBox.AvailableAfter(createdAt)); err != nil {
		return fmt.Errorf("show new box info: %w", err)
	}

	return nil
}

func calculateWaitPeriod(standartPeriod time.Duration, isImmediate bool) time.Duration {
	if isImmediate {
		return 0
	}

	return standartPeriod
}

func createNormalBox(
	chatID msginfo.ChatID,
	createdAt time.Time,
	commonWaitPeriod time.Duration,
) box.Box {
	var (
		availableAt         = createdAt
		readyNotificationAt time.Time
	)

	if commonWaitPeriod == 0 {
		readyNotificationAt = createdAt
	} else {
		availableAt = availableAt.Add(commonWaitPeriod)
	}

	return box.Box{
		ChatID:              chatID,
		Status:              box.StatusInProgress,
		Type:                box.TypeNormal,
		CreatedAt:           createdAt,
		ReadyNotificationAt: readyNotificationAt,
		AvailableAt:         availableAt,
	}
}
