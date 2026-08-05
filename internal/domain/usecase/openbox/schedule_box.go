package openbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (o *OpenBox) ScheduleBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
) error {
	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		inProgressBoxes, err := o.inProgressBoxesByType(ctx, chatID, boxType)
		if err != nil {
			return fmt.Errorf("get in progress boxes by type %s: %w", boxType.String(), err)
		}

		now := o.timeProvider.Now()

		if len(inProgressBoxes) != 0 {
			if err := o.notifier.ShowBoxInfo(ctx, inProgressBoxes[0], inProgressBoxes[0].AvailableAfter(now)); err != nil {
				return fmt.Errorf("show existing box info: %w", err)
			}

			return nil
		}

		if err := o.scheduleBox(
			ctx,
			chatID,
			boxType,
			now,
			false,
		); err != nil {
			return fmt.Errorf("schedule box: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (o *OpenBox) inProgressBoxesByType(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
) ([]box.Box, error) {
	boxes, err := o.repo.GetBoxesByStatus(ctx, chatID, box.StatusInProgress)
	if err != nil {
		return nil, fmt.Errorf("get active boxes: %w", err)
	}

	return filterBoxes(boxes, boxType), nil
}

func filterBoxes(boxes []box.Box, boxType box.Type) []box.Box {
	var filtered []box.Box

	for _, b := range boxes {
		if b.Type == boxType {
			filtered = append(filtered, b)
		}
	}

	return filtered
}

func (o *OpenBox) scheduleBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	createdAt time.Time,
	isImmediate bool,
) error {
	newBox := o.createBox(
		chatID,
		createdAt,
		boxType,
		isImmediate,
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

func (o *OpenBox) createBox(
	chatID msginfo.ChatID,
	createdAt time.Time,
	boxType box.Type,
	isImmediate bool,
) box.Box {
	var (
		availableAt         = createdAt
		readyNotificationAt time.Time
	)

	if isImmediate {
		readyNotificationAt = createdAt
	} else {
		availableAt = availableAt.Add(o.boxWaitPeriod[boxType])
	}

	return box.Box{
		ChatID:              chatID,
		Status:              box.StatusInProgress,
		Type:                boxType,
		CreatedAt:           createdAt,
		ReadyNotificationAt: readyNotificationAt,
		AvailableAt:         availableAt,
	}
}
