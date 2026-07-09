package openbox

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (o *OpenBox) Open(
	ctx context.Context,
	id int,
	chatID msginfo.ChatID,
) error {
	readyBox, err := o.repo.GetBoxByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get box by id: %w", err)
	}

	switch readyBox.Status {
	case box.StatusPending:
		return perror.InvalidStatus("your box is in pending status")

	case box.StatusOpened:
		return perror.InvalidStatus("box already opened")

	case box.StatusCanceled:
		return perror.InvalidStatus("box already canceled")

	case box.StatusInProgress:
	}

	now := o.timeProvider.Now()

	if readyBox.AvailableAt.After(now) {
		if err := o.notifier.ShowBoxInfo(ctx, readyBox, readyBox.AvailableAt.Sub(now)); err != nil {
			return fmt.Errorf("show box info: %w", err)
		}

		return nil
	}

	readyBox.Status = box.StatusOpened
	readyBox.CompletedAt = now

	if err := o.repo.UpdateBox(ctx, readyBox); err != nil {
		return fmt.Errorf("update box: %w", err)
	}

	return nil
}
