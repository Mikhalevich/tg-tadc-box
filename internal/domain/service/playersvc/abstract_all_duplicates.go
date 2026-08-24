package playersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (s *Service) AbstractAllDuplicates(
	ctx context.Context,
	chatID msginfo.ChatID,
) (gloink.Amount, player.Wallet, error) {
	var (
		abstractedGloinksAmount = gloink.AmountFromInt(0)
		wallet                  player.Wallet
	)

	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		profile, err := s.repo.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get player by chat id: %w", err)
		}

		abstractedGloinksAmount = profile.Profile.AbstractDuplicatesAll(s.abstractionCosts)

		if err := s.repo.UpdatePlayer(ctx, profile); err != nil {
			return fmt.Errorf("update player: %w", err)
		}

		wallet = profile.Profile.Wallet

		return nil
	}); err != nil {
		return 0, player.Wallet{}, fmt.Errorf("transaction: %w", err)
	}

	return abstractedGloinksAmount, wallet, nil
}
