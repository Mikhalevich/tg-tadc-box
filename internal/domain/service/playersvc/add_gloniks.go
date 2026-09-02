package playersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Service) AddGloinks(
	ctx context.Context,
	chatID msginfo.ChatID,
	amount gloink.Amount,
) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		profile, err := s.repo.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get player by chat id: %w", err)
		}

		profile.Profile.Wallet.IncreaseGloinks(amount)

		if err := s.repo.UpdatePlayer(ctx, profile); err != nil {
			return fmt.Errorf("update player: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
