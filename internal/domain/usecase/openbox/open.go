package openbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

const (
	boxWaitPeriod = 3 * time.Hour
)

func (o *OpenBox) Open(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	boxes, err := o.repo.GetBoxesByStatus(ctx, chatID, box.StatusInProgress)
	if err != nil {
		return fmt.Errorf("get active boxes: %w", err)
	}

	now := o.timeProvider.Now()

	if len(boxes) > 0 {
		if err := o.notifier.ShowBoxInfo(ctx, boxes[0], boxes[0].AvailableAfter(now)); err != nil {
			return fmt.Errorf("show existing box info: %w", err)
		}

		return nil
	}

	newBox := createNormalBox(chatID, now)

	boxID, err := o.repo.InsertBox(ctx, newBox)
	if err != nil {
		return fmt.Errorf("insert box: %w", err)
	}

	newBox.ID = boxID

	if err := o.notifier.ShowBoxInfo(ctx, newBox, newBox.AvailableAfter(now)); err != nil {
		return fmt.Errorf("show new box info: %w", err)
	}

	return nil
}

func createNormalBox(
	chatID msginfo.ChatID,
	createdAt time.Time,
) box.Box {
	return box.Box{
		ChatID:      chatID,
		Status:      box.StatusInProgress,
		Type:        box.TypeNormal,
		CreatedAt:   createdAt,
		AvailableAt: createdAt.Add(boxWaitPeriod),
	}
}
