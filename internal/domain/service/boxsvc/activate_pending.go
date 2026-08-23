package boxsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Service) ActivatePending(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
	now time.Time,
) error {
	if err := s.repo.ChangeFirstBoxStatusByType(
		ctx,
		chatID,
		boxType,
		box.StatusInProgress,
		box.StatusPending,
		now.Add(s.boxWaitPeriod[boxType]),
	); err != nil {
		return fmt.Errorf("change first box status by type: %w", err)
	}

	return nil
}
