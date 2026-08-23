package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

func (s *Service) ShowReadyToOpenBox(
	ctx context.Context,
	domBox box.Box,
) error {
	if err := s.sendBoxIsAvailable(ctx, domBox); err != nil {
		return fmt.Errorf("send box is available: %w", err)
	}

	return nil
}
