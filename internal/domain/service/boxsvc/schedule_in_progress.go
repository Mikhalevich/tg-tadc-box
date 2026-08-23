package boxsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

// ScheduleInProgress create a new box in InProgress status
// returns created box and error.
func (s *Service) ScheduleInProgress(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	createdAt time.Time,
	isImmediate bool,
) (box.Box, error) {
	var (
		newBox box.Box
		err    error
	)

	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		newBox, err = s.processScheduleInProgressBox(
			ctx,
			chatID,
			boxType,
			createdAt,
			isImmediate,
		)

		if err != nil {
			return fmt.Errorf("process schedule box: %w", err)
		}

		return nil
	}); err != nil {
		return box.Box{}, fmt.Errorf("transaction: %w", err)
	}

	return newBox, nil
}

// processScheduleBox schedule box
// returns is sheduled flag and error.
func (s *Service) processScheduleInProgressBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	createdAt time.Time,
	isImmediate bool,
) (box.Box, error) {
	inProgressBoxes, err := s.repo.GetBoxesByStatus(ctx, chatID, box.StatusInProgress)
	if err != nil {
		return box.Box{}, fmt.Errorf("get in_progress boxes: %w", err)
	}

	if len(filterBoxes(inProgressBoxes, boxType)) > 0 {
		return box.Box{}, nil
	}

	newBox, err := s.scheduleBox(
		ctx,
		chatID,
		boxType,
		box.StatusInProgress,
		createdAt,
		box.Meta{},
		isImmediate,
	)

	if err != nil {
		return box.Box{}, fmt.Errorf("schedule box: %w", err)
	}

	return newBox, nil
}
