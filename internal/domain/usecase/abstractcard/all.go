package abstractcard

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (ac *AbstractCard) All(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := ac.transactor.Transaction(ctx, func(ctx context.Context) error {
		abstractedGloinksAmount, wallet, err := ac.playerService.AbstractAllDuplicates(
			ctx,
			chatID,
		)

		if err != nil {
			return fmt.Errorf("abstract all duplicates: %w", err)
		}

		if err := ac.notifier.ShowGloinksWalletAfterAbstraction(
			ctx,
			chatID,
			wallet,
			abstractedGloinksAmount,
		); err != nil {
			return fmt.Errorf("show gloinks wallet after abstraction: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
