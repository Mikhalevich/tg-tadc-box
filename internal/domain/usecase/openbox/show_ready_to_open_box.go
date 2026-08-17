package openbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (o *OpenBox) ShowReadyToOpenBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxID box.ID,
) error {
	readyBox, err := o.repo.GetBoxByID(ctx, boxID)
	if err != nil {
		return fmt.Errorf("get box by id: %w", err)
	}

	if err := isBoxReadyToOpen(readyBox, o.timeProvider.Now()); err != nil {
		return fmt.Errorf("check box is ready for open: %w", err)
	}

	if err := o.notifier.ShowReadyToOpenBox(ctx, readyBox); err != nil {
		return fmt.Errorf("show ready to open notification: %w", err)
	}

	return nil
}

func isBoxReadyToOpen(readyBox box.Box, now time.Time) error {
	if err := isInProgress(readyBox.Status); err != nil {
		return fmt.Errorf("not is in_progress status: %w", err)
	}

	if readyBox.AvailableAt.After(now) {
		return perror.InvalidParam("box is not ready yet")
	}

	return nil
}
