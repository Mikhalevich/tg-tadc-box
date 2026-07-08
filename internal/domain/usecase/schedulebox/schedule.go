package schedulebox

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

const (
	boxWaitPeriod = 3 * time.Hour
)

func (s *ScheduleBox) Schedule(chatID msginfo.ChatID) error {
	boxes, err := s.repo.GetActiveBoxes(chatID)
	if err != nil {
		return fmt.Errorf("get active boxes: %w", err)
	}

	now := s.timeProvider.Now()

	if len(boxes) > 0 {
		s.notifier.ShowBoxInfo(boxes[0], boxes[0].AvailableAt.Sub(now))

		return nil
	}

	newBox := createNormalBox(chatID, now)

	boxID, err := s.repo.InsertBox(newBox)
	if err != nil {
		return fmt.Errorf("insert box: %w", err)
	}

	newBox.ID = boxID

	s.notifier.ShowBoxInfo(newBox, newBox.AvailableAt.Sub(now))

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
