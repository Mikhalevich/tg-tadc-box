package outboxsvc

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/outboxmsg"
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type Repository interface {
	OutboxSelectForDispatchMessages(
		ctx context.Context,
		visibilityAt time.Time,
		limit int,
	) ([]outboxmsg.Message, error)

	OutboxUpdateStatus(
		ctx context.Context,
		ids []int,
		status outboxmsg.Status,
		updatedAt time.Time,
	) error

	OutboxIncrementRetryCount(
		ctx context.Context,
		ids []int,
		updatedAt time.Time,
	) error
}

type Sender interface {
	SendMessage(ctx context.Context, msg msginfo.Message) error
}

type TimeProvider interface {
	Now() time.Time
}

type Service struct {
	transactor   Transactor
	repo         Repository
	sender       Sender
	timeProvider TimeProvider
}

func New(
	transactor Transactor,
	repo Repository,
	sender Sender,
	timeProvider TimeProvider,
) *Service {
	return &Service{
		transactor:   transactor,
		repo:         repo,
		sender:       sender,
		timeProvider: timeProvider,
	}
}
