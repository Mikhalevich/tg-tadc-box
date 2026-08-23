package playersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (s *Service) UpdatePlayer(
	ctx context.Context,
	usr player.Player,
) error {
	if err := s.repo.UpdatePlayer(ctx, usr); err != nil {
		return fmt.Errorf("repo update player: %w", err)
	}

	return nil
}
