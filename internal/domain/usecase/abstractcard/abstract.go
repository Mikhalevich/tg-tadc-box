package abstractcard

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (ac *AbstractCard) Abstract(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	rewardType reward.RewardType,
	pos int,
	count int,
) error {
	if err := ac.transactor.Transaction(ctx, func(ctx context.Context) error {
		abstractedGloinksAmount, wallet, err := ac.playerService.AbstractCard(
			ctx,
			chatID,
			rewardType,
			pos,
			count,
		)

		if err != nil {
			return fmt.Errorf("abstract card: %w", err)
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

	if err := ac.pageProvider.Page(
		ctx,
		chatID,
		messageID,
		rewardType,
		pos,
	); err != nil {
		return fmt.Errorf("show page: %w", err)
	}

	return nil
}
