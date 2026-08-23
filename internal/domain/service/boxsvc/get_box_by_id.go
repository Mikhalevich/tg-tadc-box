package boxsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

func (s *Service) GetBoxByID(
	ctx context.Context,
	id box.ID,
) (box.Box, error) {
	boxFromRepo, err := s.repo.GetBoxByID(ctx, id)
	if err != nil {
		return box.Box{}, fmt.Errorf("get box by id: %w", err)
	}

	return boxFromRepo, nil
}
