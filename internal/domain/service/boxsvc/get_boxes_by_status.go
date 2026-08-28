package boxsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Service) GetBoxesByStatus(
	ctx context.Context,
	chatID msginfo.ChatID,
	statuses ...box.Status,
) ([]box.Box, error) {
	boxes, err := s.repo.GetBoxesByStatus(
		ctx,
		chatID,
		statuses...,
	)

	if err != nil {
		return nil, fmt.Errorf("get boxes from repo: %w", err)
	}

	return boxes, nil
}
