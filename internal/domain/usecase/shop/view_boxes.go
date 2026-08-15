package shop

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Shop) ViewBoxes(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	plr, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	inProgressBoxes, err := s.boxProvider.GetBoxesByStatus(
		ctx,
		chatID,
		box.StatusInProgress,
	)

	if err != nil {
		return fmt.Errorf("get in_progress boxes: %w", err)
	}

	if err := s.notifier.ShowBoxCosts(
		ctx,
		chatID,
		messageID,
		plr.Profile.Wallet,
		s.boxCosts,
		box.ToInProgressBoxMapByType(inProgressBoxes, s.timeProvider.Now()),
	); err != nil {
		return fmt.Errorf("show box costs: %w", err)
	}

	return nil
}
