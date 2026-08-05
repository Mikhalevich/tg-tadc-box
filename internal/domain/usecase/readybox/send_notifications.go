package readybox

import (
	"context"
	"fmt"
	"time"
)

func (r *ReadyBox) SendNotifications(
	ctx context.Context,
	limit int,
) error {
	if err := r.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := r.processTransaction(
			ctx,
			r.timeProvider.Now(),
			limit,
		); err != nil {
			return fmt.Errorf("process ready to open boxes: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (r *ReadyBox) processTransaction(
	ctx context.Context,
	now time.Time,
	limit int,
) error {
	readyBoxes, err := r.repo.GetReadyToOpenBoxes(ctx, r.timeProvider.Now(), limit)
	if err != nil {
		return fmt.Errorf("get ready to open boxes: %w", err)
	}

	if len(readyBoxes) == 0 {
		return nil
	}

	ids := make([]int, 0, len(readyBoxes))

	for _, readyBox := range readyBoxes {
		if err := r.notifier.ShowReadyToOpenBox(ctx, readyBox); err != nil {
			return fmt.Errorf("show box info: %w", err)
		}

		ids = append(ids, readyBox.ID.Int())
	}

	if err := r.repo.SetBoxReadyNotificationAt(ctx, ids, now); err != nil {
		return fmt.Errorf("set box notification at: %w", err)
	}

	return nil
}
