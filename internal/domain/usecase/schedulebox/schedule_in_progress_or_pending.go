package schedulebox

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

// ScheduleInProgressOrPending create a new box in InProgress or Pending status
// returns created box and error.
func (s *ScheduleBox) ScheduleInProgressOrPending(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	meta box.Meta,
) (box.Box, error) {
	var (
		newBox box.Box
		err    error
	)
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		newBox, err = s.processScheduleInProgressOrPendingBox(
			ctx,
			chatID,
			boxType,
			meta,
		)

		if err != nil {
			return fmt.Errorf("process schedule in progress or pending box: %w", err)
		}

		return nil
	}); err != nil {
		return box.Box{}, fmt.Errorf("transaction: %w", err)
	}

	return newBox, nil
}

func (s *ScheduleBox) processScheduleInProgressOrPendingBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	meta box.Meta,
) (box.Box, error) {
	inProgressBoxes, err := s.repo.GetBoxesByStatus(ctx, chatID, box.StatusInProgress)
	if err != nil {
		return box.Box{}, fmt.Errorf("get in_progress boxes: %w", err)
	}

	boxStatus := box.StatusInProgress

	if len(filterBoxes(inProgressBoxes, boxType)) > 0 {
		boxStatus = box.StatusPending
	}

	newBox, err := s.scheduleBox(
		ctx,
		chatID,
		boxType,
		boxStatus,
		s.timeProvider.Now(),
		meta,
		false,
	)
	if err != nil {
		return box.Box{}, fmt.Errorf("schedule box: %w", err)
	}

	return newBox, nil
}
