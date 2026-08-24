package likesvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/like"
)

func (s *Service) AddLike(
	ctx context.Context,
	rwdLike like.Like,
) error {
	if err := s.repo.InsertLike(ctx, rwdLike); err != nil {
		return fmt.Errorf("insert like: %w", err)
	}

	return nil
}
