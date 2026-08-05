package shop

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Shop) ViewBoxes(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	plr, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	if err := s.notifier.ShowBoxCosts(
		ctx,
		chatID,
		plr.Profile.Wallet,
		s.boxCosts,
	); err != nil {
		return fmt.Errorf("show box costs: %w", err)
	}

	return nil
}
