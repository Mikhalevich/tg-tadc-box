package playerprovider

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

// GetPlayerByChatID get existing or create a new player
// returns player and error.
func (p *PlayerProvider) GetPlayerByChatID(
	ctx context.Context,
	chatID msginfo.ChatID,
) (player.Player, error) {
	plr, err := p.repo.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return player.Player{}, fmt.Errorf("get player by chat id: %w", err)
		}

		plr, err := p.createPlayer(ctx, chatID, p.timeProvider.Now())
		if err != nil {
			return player.Player{}, fmt.Errorf("create player: %w", err)
		}

		return plr, nil
	}

	return plr, nil
}

func (p *PlayerProvider) createPlayer(
	ctx context.Context,
	chatID msginfo.ChatID,
	createdAt time.Time,
) (player.Player, error) {
	plr := player.Player{
		ChatID:           chatID,
		CreatedAt:        createdAt,
		ProfileUpdatedAt: createdAt,
	}

	id, err := p.repo.InsertPlayer(ctx, plr)
	if err != nil {
		return player.Player{}, fmt.Errorf("insert player: %w", err)
	}

	plr.ID = player.IDFromInt(id)

	return plr, nil
}
