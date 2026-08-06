package shop

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (s *Shop) BuyBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		plr, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get player by chat_id: %w", err)
		}

		amount, err := s.findBoxCost(boxType)
		if err != nil {
			return fmt.Errorf("find box cost: %w", err)
		}

		if !amount.IsFree() {
			if err := plr.Profile.Wallet.DecreaseGloinks(amount); err != nil {
				return fmt.Errorf("decrease gloinks: %w", err)
			}
		}

		isScheduled, err := s.boxScheduler.ScheduleBox(
			ctx,
			chatID,
			boxType,
			plr.Profile.IsNoOpenBoxes(),
		)
		if err != nil {
			return fmt.Errorf("schedule box: %w", err)
		}

		if !isScheduled || amount.IsFree() {
			return nil
		}

		if err := s.playerProvider.UpdatePlayer(ctx, plr); err != nil {
			return fmt.Errorf("update player: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (s *Shop) findBoxCost(boxType box.Type) (gloink.Amount, error) {
	for _, cost := range s.boxCosts {
		if cost.Type == boxType {
			return cost.Amount, nil
		}
	}

	return gloink.Amount(0), perror.InvalidParam("invalid box type")
}
