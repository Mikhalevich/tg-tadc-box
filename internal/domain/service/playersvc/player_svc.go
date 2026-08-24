package playersvc

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

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

type Service struct {
	abstractionCosts map[reward.RewardType]gloink.Amount
	transactor       Transactor
	repo             Repository
	timeProvider     TimeProvider
}

func New(
	abstractionCosts map[reward.RewardType]gloink.Amount,
	transactor Transactor,
	repo Repository,
	timeProvider TimeProvider,
) *Service {
	return &Service{
		abstractionCosts: abstractionCosts,
		transactor:       transactor,
		repo:             repo,
		timeProvider:     timeProvider,
	}
}
