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
		profile, err := ac.playerProvider.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get player by chat id: %w", err)
		}

		abstractedGloinksAmount := profile.Profile.AbstractDuplicatesAll()

		if err := ac.playerProvider.UpdatePlayer(ctx, profile); err != nil {
			return fmt.Errorf("update player: %w", err)
		}

		if err := ac.notifier.ShowGloinksWalletAfterAbstraction(
			ctx,
			chatID,
			profile.Profile.Wallet,
			abstractedGloinksAmount,
		); err != nil {
			return fmt.Errorf("show gloiks wallet after abstraction")
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
