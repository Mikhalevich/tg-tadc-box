package boxsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

// ScheduleInProgressOrPending create a new box in InProgress or Pending status
// returns created box and error.
func (s *Service) ScheduleInProgressOrPending(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	createdAt time.Time,
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
			createdAt,
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

func (s *Service) processScheduleInProgressOrPendingBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	createdAt time.Time,
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
		createdAt,
		meta,
		false,
	)
	if err != nil {
		return box.Box{}, fmt.Errorf("schedule box: %w", err)
	}

	return newBox, nil
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

func (s *Service) scheduleBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	boxStatus box.Status,
	createdAt time.Time,
	meta box.Meta,
	isImmediate bool,
) (box.Box, error) {
	newBox := s.createBox(
		chatID,
		createdAt,
		boxType,
		boxStatus,
		meta,
		isImmediate,
	)

	boxID, err := s.repo.InsertBox(ctx, newBox)
	if err != nil {
		return box.Box{}, fmt.Errorf("insert box: %w", err)
	}

	newBox.ID = box.IDFromInt(boxID)

	return newBox, nil
}

func (s *Service) createBox(
	chatID msginfo.ChatID,
	createdAt time.Time,
	boxType box.Type,
	boxStatus box.Status,
	meta box.Meta,
	isImmediate bool,
) box.Box {
	var (
		availableAt         = createdAt
		readyNotificationAt time.Time
	)

	if isImmediate {
		readyNotificationAt = createdAt
	} else {
		availableAt = availableAt.Add(s.boxWaitPeriod[boxType])
	}

	return box.Box{
		ChatID:              chatID,
		Status:              boxStatus,
		Type:                boxType,
		CreatedAt:           createdAt,
		ReadyNotificationAt: readyNotificationAt,
		AvailableAt:         availableAt,
		Meta:                meta,
	}
}
