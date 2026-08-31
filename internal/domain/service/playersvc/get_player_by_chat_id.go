package playersvc

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

const (
	initialGloinksAmount = 100
)

// GetPlayerByChatID get existing or create a new player
// returns player and error.
func (s *Service) GetPlayerByChatID(
	ctx context.Context,
	chatID msginfo.ChatID,
) (player.Player, error) {
	plr, err := s.getOrCreatePlayer(ctx, chatID)
	if err != nil {
		return player.Player{}, fmt.Errorf("get or create player: %w", err)
	}

	return plr, nil
}

func (s *Service) getOrCreatePlayer(
	ctx context.Context,
	chatID msginfo.ChatID,
) (player.Player, error) {
	plr, err := s.repo.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return player.Player{}, fmt.Errorf("get player by chat id: %w", err)
		}

		plr, err := s.createPlayer(ctx, chatID, s.timeProvider.Now())
		if err != nil {
			return player.Player{}, fmt.Errorf("create player: %w", err)
		}

		return plr, nil
	}

	return plr, nil
}

func (s *Service) createPlayer(
	ctx context.Context,
	chatID msginfo.ChatID,
	createdAt time.Time,
) (player.Player, error) {
	plr := player.Player{
		ChatID:    chatID,
		CreatedAt: createdAt,
		Profile: player.Profile{
			Wallet: player.Wallet{
				GloinksAmount: gloink.AmountFromInt(initialGloinksAmount),
			},
		},
		ProfileUpdatedAt: createdAt,
	}

	id, err := s.repo.InsertPlayer(ctx, plr)
	if err != nil {
		return player.Player{}, fmt.Errorf("insert player: %w", err)
	}

	plr.ID = player.IDFromInt(id)

	return plr, nil
}
