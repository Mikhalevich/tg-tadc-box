package playerprovider

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

type Repository interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
	InsertPlayer(ctx context.Context, usr player.Player) (int, error)
	UpdatePlayer(ctx context.Context, usr player.Player) error
}

type TimeProvider interface {
	Now() time.Time
}

type PlayerProvider struct {
	repo         Repository
	timeProvider TimeProvider
}

func New(
	repo Repository,
	timeProvider TimeProvider,
) *PlayerProvider {
	return &PlayerProvider{
		repo:         repo,
		timeProvider: timeProvider,
	}
}
