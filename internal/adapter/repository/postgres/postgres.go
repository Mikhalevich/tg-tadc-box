package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/boxsvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/likesvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/messagesvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/outbox/outboxsvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/playersvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/referralsvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/rewardsvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/readybox"
)

var (
	_ rewardsvc.Repository        = (*Postgres)(nil)
	_ playersvc.Repository        = (*Postgres)(nil)
	_ boxsvc.Repository           = (*Postgres)(nil)
	_ messagesvc.ButtonRepository = (*Postgres)(nil)
	_ outboxsvc.Repository        = (*Postgres)(nil)
	_ likesvc.Repository          = (*Postgres)(nil)
	_ referralsvc.Repository      = (*Postgres)(nil)

	_ readybox.Repository = (*Postgres)(nil)
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

func (p *Postgres) Ping(ctx context.Context) error {
	var (
		query = `
			SELECT
				id
			FROM
				player
			LIMIT
				1
			
		`

		res []int
	)
	if err := sqlx.SelectContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&res,
		query,
	); err != nil {
		return fmt.Errorf("select player table: %w", err)
	}

	return nil
}
