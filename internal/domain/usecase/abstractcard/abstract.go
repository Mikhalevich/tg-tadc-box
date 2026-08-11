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
		profile, err := ac.playerProvider.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get player by chat id: %w", err)
		}

		abstractedGloinksAmount, err := profile.Profile.AbstractByPos(
			rewardType,
			pos,
			count,
			ac.abstractionCosts,
		)

		if err != nil {
			return fmt.Errorf("abstract by pos: %w", err)
		}

		if err := ac.playerProvider.UpdatePlayer(ctx, profile); err != nil {
			return fmt.Errorf("update player: %w", err)
		}

		if err := ac.notifier.ShowGloinksWalletAfterAbstraction(
			ctx,
			chatID,
			profile.Profile.Wallet,
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
