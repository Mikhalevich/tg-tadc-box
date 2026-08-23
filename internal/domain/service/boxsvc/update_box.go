package boxsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

func (s *Service) UpdateBox(
	ctx context.Context,
	b box.Box,
) error {
	if err := s.repo.UpdateBox(ctx, b); err != nil {
		return fmt.Errorf("update box: %w", err)
	}

	return nil
}
