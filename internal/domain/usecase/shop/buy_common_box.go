package shop

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Shop) BuyCommonBox(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		plr, _, err := s.playerService.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get player by chat_id: %w", err)
		}

		isImmediate := isImmediateOpen(box.TypeCommon, plr)

		scheduledBox, err := s.boxService.ScheduleInProgress(
			ctx,
			chatID,
			box.TypeCommon,
			s.timeProvider.Now(),
			isImmediate,
		)

		if err != nil {
			return fmt.Errorf("schedule box: %w", err)
		}

		if err := s.showBuyBoxNotification(
			ctx,
			plr,
			scheduledBox,
			isImmediate,
		); err != nil {
			return fmt.Errorf("show buy box notification: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
