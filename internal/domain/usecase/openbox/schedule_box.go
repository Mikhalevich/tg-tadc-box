package openbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

// ScheduleBox schedule box
// returns is sheduled flag and error.
func (o *OpenBox) ScheduleBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	isImmediate bool,
) (bool, error) {
	var (
		isScheduled = false
		err         error
	)

	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		isScheduled, err = o.processScheduleBox(
			ctx,
			chatID,
			boxType,
			isImmediate,
		)

		if err != nil {
			return fmt.Errorf("process schedule box: %w", err)
		}

		return nil
	}); err != nil {
		return false, fmt.Errorf("transaction: %w", err)
	}

	return isScheduled, nil
}

// processScheduleBox schedule box
// returns is sheduled flag and error.
func (o *OpenBox) processScheduleBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	isImmediate bool,
) (bool, error) {
	inProgressBoxes, err := o.repo.GetBoxesByStatus(ctx, chatID, box.StatusInProgress)
	if err != nil {
		return false, fmt.Errorf("get in_progress boxes: %w", err)
	}

	if len(filterBoxes(inProgressBoxes, boxType)) != 0 {
		if err := o.showInProgressBoxesNotification(
			ctx,
			chatID,
			inProgressBoxes,
		); err != nil {
			return false, fmt.Errorf("show in_progress notification: %w", err)
		}

		return false, nil
	}

	newBox, err := o.scheduleBox(
		ctx,
		chatID,
		boxType,
		o.timeProvider.Now(),
		isImmediate,
	)
	if err != nil {
		return false, fmt.Errorf("schedule box: %w", err)
	}

	inProgressBoxes = append(inProgressBoxes, newBox)

	if err := o.showInProgressBoxesNotification(
		ctx,
		chatID,
		inProgressBoxes,
	); err != nil {
		return false, fmt.Errorf("show in_progress notification: %w", err)
	}

	return true, nil
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
) (box.Box, error) {
	newBox := o.createBox(
		chatID,
		createdAt,
		boxType,
		isImmediate,
	)

	boxID, err := o.repo.InsertBox(ctx, newBox)
	if err != nil {
		return box.Box{}, fmt.Errorf("insert box: %w", err)
	}

	newBox.ID = box.IDFromInt(boxID)

	return newBox, nil
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
