package shop

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (s *Shop) ViewBoxes(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	plr, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	if err := s.showPlayerBoxes(
		ctx,
		plr,
	); err != nil {
		return fmt.Errorf("show boxes: %w", err)
	}

	return nil
}

func (s *Shop) showPlayerBoxes(
	ctx context.Context,
	plr player.Player,
) error {
	inProgressBoxes, err := s.boxProvider.GetBoxesByStatus(
		ctx,
		plr.ChatID,
		box.StatusInProgress,
	)

	if err != nil {
		return fmt.Errorf("get in_progress boxes: %w", err)
	}

	if err := s.notifier.ShowBoxCosts(
		ctx,
		plr.ChatID,
		plr.Profile.Wallet,
		s.boxCosts,
		box.ToInProgressBoxMapByType(inProgressBoxes, s.timeProvider.Now()),
	); err != nil {
		return fmt.Errorf("show box costs: %w", err)
	}

	return nil
}
