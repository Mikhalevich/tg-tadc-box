package openbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (o *OpenBox) ShowInProgressBoxes(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	boxes, err := o.repo.GetBoxesByStatus(ctx, chatID, box.StatusInProgress)
	if err != nil {
		return fmt.Errorf("get in progress boxes: %w", err)
	}

	if err := o.showInProgressBoxesNotification(
		ctx,
		chatID,
		boxes,
	); err != nil {
		return fmt.Errorf("show in progress boxes: %w", err)
	}

	return nil
}

func (o *OpenBox) showInProgressBoxesNotification(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxes []box.Box,
) error {
	box.SortBoxByType(boxes)

	if err := o.notifier.ShowInProgressBoxes(
		ctx,
		chatID,
		convertToInProgressBoxes(boxes, o.timeProvider.Now()),
	); err != nil {
		return fmt.Errorf("show in progress boxes: %w", err)
	}

	return nil
}

func convertToInProgressBoxes(boxes []box.Box, now time.Time) []box.InProgressBox {
	inProgressBoxes := make([]box.InProgressBox, 0, len(boxes))

	for _, b := range boxes {
		inProgressBoxes = append(inProgressBoxes, box.InProgressBox{
			Box:            b,
			AvailableAfter: b.AvailableAfter(now),
		})
	}

	return inProgressBoxes
}
