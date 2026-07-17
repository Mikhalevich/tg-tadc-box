package setup

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-tadc-box/cmd/poller/internal/app"
	"github.com/Mikhalevich/tg-tadc-box/cmd/poller/internal/config"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/imageprovider"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/notifier"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/readybox"
)

func StartWorker(ctx context.Context, cfg config.Config) error {
	botAPI, err := bot.New(cfg.Bot.Token, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot api: %w", err)
	}

	pgDB, dbCleanup, err := MakePostgres(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}

	defer dbCleanup()

	var (
		messageSender       = messagesender.New(botAPI)
		messageProcessor    = messageprocessor.New(messageSender, messageSender, pgDB, nil)
		imageProvider       = imageprovider.New()
		notificationService = notifier.New(messageProcessor, messageProcessor, imageProvider)
		timeProvider        = timeprovider.New()
		readyBoxService     = readybox.New(pgDB, pgDB.Transactor(), notificationService, timeProvider)
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
