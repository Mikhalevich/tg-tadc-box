package shop

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (s *Shop) BuyBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	rewardType reward.RewardType,
) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		plr, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get player by chat_id: %w", err)
		}

		amount, err := s.findBoxCost(rewardType)
		if err != nil {
			return fmt.Errorf("find box cost: %w", err)
		}

		if err := plr.Profile.Wallet.DecreaseGloinks(amount); err != nil {
			return fmt.Errorf("decrease gloinks: %w", err)
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

func (s *Shop) findBoxCost(rewardType reward.RewardType) (gloink.Amount, error) {
	for _, cost := range s.boxCosts {
		if cost.Type == rewardType {
			return cost.Amount, nil
		}
	}

	return gloink.Amount(0), perror.InvalidParam("invalid box type")
}
