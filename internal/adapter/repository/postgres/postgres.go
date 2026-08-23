package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/rewardsvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/likereward"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/readybox"
)

var (
	_ readybox.Repository   = (*Postgres)(nil)
	_ rewardsvc.Repository  = (*Postgres)(nil)
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
