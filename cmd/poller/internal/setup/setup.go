package setup

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-tadc-box/cmd/poller/internal/app"
	"github.com/Mikhalevich/tg-tadc-box/cmd/poller/internal/config"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/markdownescaper"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/notificationsvc"
	outboximageprovider "github.com/Mikhalevich/tg-tadc-box/internal/domain/service/outbox/imageprovider"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/readybox"
)

func StartWorker(ctx context.Context, cfg config.Config) error {
	pgDB, dbCleanup, err := MakePostgres(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}

	defer dbCleanup()

	var (
		notificationService = notificationsvc.New(
			pgDB,
			markdownescaper.New(),
			outboximageprovider.New(),
		)
		readyBoxService = readybox.New(
			pgDB,
			pgDB.Transactor(),
			notificationService,
			timeprovider.New(),
		)
	)

	app.New(
		readyBoxService,
	).Run(
		ctx,
		cfg.BoxReadyNotificationWorker,
	)

	return nil
}

func MakePostgres(cfg config.Postgres) (*postgres.Postgres, func(), error) {
	driver := driver.NewPgx()

	dbConn, err := otelsql.Open(driver.Name(), cfg.Connection)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	var (
		sqlxDBConn = sqlx.NewDb(dbConn, driver.Name())
		transactor = transaction.New(transaction.NewSqlxDB(sqlxDBConn))
		p          = postgres.New(driver, transactor)
	)

	return p, func() {
		dbConn.Close()
	}, nil
}
