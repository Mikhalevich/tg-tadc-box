package setup

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app"
	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/config"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/imageprovider"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/markdownescaper"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/transaction"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/timeprovider"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/messageprocessor"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/notifier"
	outboximageprovider "github.com/Mikhalevich/tg-tadc-box/internal/domain/outbox/imageprovider"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/rewardgenerator"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/openbox"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/viewcards"
)

func StartBot(ctx context.Context, cfg config.Config) error {
	botAPI, err := bot.New(cfg.Bot.Token, bot.WithSkipGetMe())
	if err != nil {
		return fmt.Errorf("creating bot api: %w", err)
	}

	pgDB, dbCleanup, err := MakePostgres(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("make postgres: %w", err)
	}
	defer dbCleanup()

	boxRewardPercent, err := convertBoxRewardPercent(cfg.BoxRewardPercent)
	if err != nil {
		return fmt.Errorf("invalid box rewards percent: %w", err)
	}

	var (
		msgSender       = messagesender.New(botAPI)
		msgProcessor    = messageprocessor.New(msgSender, msgSender, pgDB)
		markdownEscaper = markdownescaper.New()
		outboxNotifier  = notifier.New(
			pgDB,
			markdownEscaper,
			outboximageprovider.New(),
		)
		boxProcessor = openbox.New(
			pgDB,
			pgDB.Transactor(),
			rewardgenerator.New(pgDB, boxRewardPercent),
			outboxNotifier,
			timeprovider.New(),
			cfg.OpenBox.CommonWaitPeriod,
		)
		directNotifier = notifier.New(
			msgProcessor,
			markdownEscaper,
			imageprovider.New(),
		)
		cardViewer = viewcards.New(
			pgDB,
			pgDB,
			directNotifier,
		)
	)

	if err := app.Start(
		ctx,
		cfg.Bot,
		msgProcessor,
		boxProcessor,
		cardViewer,
		outboxNotifier,
		outboxNotifier,
	); err != nil {
		return fmt.Errorf("app start: %w", err)
	}

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

func convertBoxRewardPercent(
	cfgRewards map[string]config.BoxRewardPercent,
) (map[box.Type]rewardgenerator.RewardPercent, error) {
	rewards := make(map[box.Type]rewardgenerator.RewardPercent, len(cfgRewards))

	for strType, percents := range cfgRewards {
		boxType, err := box.TypeFromString(strType)
		if err != nil {
			return nil, fmt.Errorf("convert box type: %w", err)
		}

		rewards[boxType] = rewardgenerator.RewardPercent{
			Legendary: percents.Legendary,
			Epic:      percents.Epic,
			Rare:      percents.Rare,
		}
	}

	return rewards, nil
}
