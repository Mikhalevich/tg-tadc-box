package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/boxsvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/messagesvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/outbox/outboxsvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/playersvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/rewardsvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/likereward"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/readybox"
)

var (
	_ rewardsvc.Repository        = (*Postgres)(nil)
	_ playersvc.Repository        = (*Postgres)(nil)
	_ boxsvc.Repository           = (*Postgres)(nil)
	_ messagesvc.ButtonRepository = (*Postgres)(nil)
	_ outboxsvc.Repository        = (*Postgres)(nil)

	_ readybox.Repository   = (*Postgres)(nil)
	_ likereward.Repository = (*Postgres)(nil)
)

type Driver interface {
	IsConstraintError(err error, constraint string) bool
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
	ExtContext(ctx context.Context) sqlx.ExtContext
}

type Postgres struct {
	dbDriver   Driver
	transactor Transactor
}

func New(
	dbDriver Driver,
	transactor Transactor,
) *Postgres {
	return &Postgres{
		transactor: transactor,
		dbDriver:   dbDriver,
	}
}

func (p *Postgres) Transactor() Transactor {
	return p.transactor
}
