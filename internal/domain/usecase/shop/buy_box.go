package shop

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (s *Shop) BuyBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.processBuyBox(
			ctx,
			chatID,
			boxType,
		); err != nil {
			return fmt.Errorf("process buy box: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (s *Shop) processBuyBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxType box.Type,
) error {
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

	isImmediate := isImmediateOpen(boxType, plr)

	scheduledBox, err := s.boxService.ScheduleInProgress(
		ctx,
		chatID,
		boxType,
		s.timeProvider.Now(),
		isImmediate,
	)
	if err != nil {
		return fmt.Errorf("schedule box: %w", err)
	}

	if scheduledBox.IsValid() && !amount.IsFree() {
		if err := s.playerProvider.UpdatePlayer(ctx, plr); err != nil {
			return fmt.Errorf("update player: %w", err)
		}
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
}

func (s *Shop) showBuyBoxNotification(
	ctx context.Context,
	plr player.Player,
	scheduledBox box.Box,
	isImmediate bool,
) error {
	if scheduledBox.IsValid() && isImmediate {
		if err := s.notifier.ShowReadyToOpenBox(
			ctx,
			scheduledBox,
		); err != nil {
			return fmt.Errorf("show ready to open box: %w", err)
		}

		return nil
	}

	if err := s.showPlayerBoxes(
		ctx,
		plr,
	); err != nil {
		return fmt.Errorf("show boxes: %w", err)
	}

	return nil
}

func isImmediateOpen(boxType box.Type, plr player.Player) bool {
	if boxType != box.TypeCommon {
		return false
	}

	return plr.Profile.IsNoOpenCommonBoxes()
}

func (s *Shop) findBoxCost(boxType box.Type) (gloink.Amount, error) {
	for _, cost := range s.boxCosts {
		if cost.Type == boxType {
			return cost.Amount, nil
		}
	}

	return gloink.Amount(0), perror.InvalidParam("invalid box type")
}
