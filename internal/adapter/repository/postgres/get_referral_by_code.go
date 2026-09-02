package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

func (p *Postgres) GetReferralByCode(
	ctx context.Context,
	code referral.Code,
) (referral.Referral, error) {
	var (
		query = `
			SELECT
				chat_id,
				invited_by_chat_id,
				code,
				created_at
			FROM
				referral
			WHERE
				code = $1
		`

		dbRef model.Referral
	)

	if err := sqlx.GetContext(
		ctx,
		p.transactor.ExtContext(ctx),
		&dbRef,
		query,
		code,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return referral.Referral{}, perror.NotFound("referral not found")
		}

		return referral.Referral{}, fmt.Errorf("get context: %w", err)
	}

	return dbRef.ToDom(), nil
}
