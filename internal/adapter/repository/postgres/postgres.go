package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/rewardgenerator"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/openbox"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/readybox"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/schedulebox"
)

var (
	_ openbox.Repository            = (*Postgres)(nil)
	_ schedulebox.Repository        = (*Postgres)(nil)
	_ readybox.Repository           = (*Postgres)(nil)
	_ rewardgenerator.RewardsGetter = (*Postgres)(nil)
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
