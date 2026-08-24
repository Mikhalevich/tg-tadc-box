package setup

import (
	"context"
	"fmt"
	"time"

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
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/boxsvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/messagesvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/notificationsvc"
	outboximageprovider "github.com/Mikhalevich/tg-tadc-box/internal/domain/service/outbox/imageprovider"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/playersvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/service/rewardsvc"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/abstractcard"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/likereward"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/openbox"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/shop"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/showreadybox"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/viewcards"
)

//nolint:funlen
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
		return fmt.Errorf("convert box rewards percent: %w", err)
	}

	boxWaitPeriod, err := convertBoxWaitPeriod(cfg.BoxWaitPeriod)
	if err != nil {
		return fmt.Errorf("convert box wait period: %w", err)
	}

	bonusBoxAttempts, err := convertBonusBoxAttempts(cfg.BonusBoxAttempts)
	if err != nil {
		return fmt.Errorf("convert box wait period: %w", err)
	}

	abstractionCosts, err := convertAbstractionCosts(cfg.AbstractionCosts)
	if err != nil {
		return fmt.Errorf("convert abstraction costs: %w", err)
	}

	var (
		msgSender                 = messagesender.New(botAPI)
		messageService            = messagesvc.New(msgSender, msgSender, pgDB)
		markdownEscaper           = markdownescaper.New()
		outboxNotificationService = notificationsvc.New(
			pgDB,
			markdownEscaper,
			outboximageprovider.New(),
		)
		notificationService = notificationsvc.New(
			messageService,
			markdownEscaper,
			imageprovider.New(),
		)
		timeProvider  = timeprovider.New()
		playerService = playersvc.New(pgDB, timeProvider)
		boxService    = boxsvc.New(
			boxWaitPeriod,
			pgDB.Transactor(),
			pgDB,
		)
		rewardService = rewardsvc.New(
			boxRewardPercent,
			pgDB,
		)
		boxOpenUsecase = openbox.New(
			bonusBoxAttempts,
			pgDB.Transactor(),
			boxService,
			playerService,
			rewardService,
			outboxNotificationService,
			timeProvider,
		)
		boxShowReadyUsecase = showreadybox.New(
			pgDB,
			notificationService,
			timeProvider,
		)
		cardViewer = viewcards.New(
			abstractionCosts,
			playerService,
			rewardService,
			notificationService,
		)
		cardAbstracter = abstractcard.New(
			abstractionCosts,
			pgDB.Transactor(),
			playerService,
			cardViewer,
			notificationService,
		)
		shopBox = shop.New(
			convertBoxCosts(cfg.BoxCosts),
			pgDB.Transactor(),
			playerService,
			boxService,
			pgDB,
			outboxNotificationService,
			timeProvider,
		)
		likeReward = likereward.New(
			pgDB.Transactor(),
			pgDB,
			outboxNotificationService,
			timeProvider,
		)
	)

	if err := app.Start(
		ctx,
		cfg.Bot,
		messageService,
		boxOpenUsecase,
		boxShowReadyUsecase,
		cardViewer,
		cardAbstracter,
		shopBox,
		likeReward,
		outboxNotificationService,
		outboxNotificationService,
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
) (map[box.Type]rewardsvc.RewardPercent, error) {
	rewards := make(map[box.Type]rewardsvc.RewardPercent, len(cfgRewards))

	for strType, percents := range cfgRewards {
		boxType, err := box.TypeFromString(strType)
		if err != nil {
			return nil, fmt.Errorf("convert box type: %w", err)
		}

		rewards[boxType] = rewardsvc.RewardPercent{
			Legendary: percents.Legendary,
			Epic:      percents.Epic,
			Rare:      percents.Rare,
		}
	}

	return rewards, nil
}

func convertBoxCosts(
	cfgBoxCosts map[string]config.BoxCost,
) []gloink.BoxCost {
	var (
		boxCosts        = make([]gloink.BoxCost, 0, len(cfgBoxCosts))
		boxTypeOrdering = []box.Type{
			box.TypeCommon,
			box.TypeRare,
			box.TypeEpic,
			box.TypeLegendary,
		}
	)

	for _, boxType := range boxTypeOrdering {
		cost, ok := cfgBoxCosts[boxType.String()]
		if !ok {
			continue
		}

		boxCosts = append(boxCosts, gloink.BoxCost{
			Type:   boxType,
			Amount: gloink.AmountFromInt(cost.Amount),
		})
	}

	return boxCosts
}

func convertBoxWaitPeriod(
	cfgBoxWaitPeriod map[string]time.Duration,
) (map[box.Type]time.Duration, error) {
	boxWaitPeriod := make(map[box.Type]time.Duration, len(cfgBoxWaitPeriod))

	for boxTypeRaw, waitPeriod := range cfgBoxWaitPeriod {
		boxType, err := box.TypeFromString(boxTypeRaw)
		if err != nil {
			return nil, fmt.Errorf("convert to box type: %w", err)
		}

		boxWaitPeriod[boxType] = waitPeriod
	}

	return boxWaitPeriod, nil
}

func convertBonusBoxAttempts(
	cfgBonusBoxAttempts map[string]int,
) (map[box.Type]int, error) {
	bonusBoxAttempts := make(map[box.Type]int, len(cfgBonusBoxAttempts))

	for boxTypeRaw, attempts := range cfgBonusBoxAttempts {
		boxType, err := box.TypeFromString(boxTypeRaw)
		if err != nil {
			return nil, fmt.Errorf("convert to box type: %w", err)
		}

		bonusBoxAttempts[boxType] = attempts
	}

	return bonusBoxAttempts, nil
}

func convertAbstractionCosts(
	cfgCosts map[string]int,
) (map[reward.RewardType]gloink.Amount, error) {
	costs := make(map[reward.RewardType]gloink.Amount, len(cfgCosts))

	for rawRewardType, rawAmount := range cfgCosts {
		rewardType, err := reward.TypeFromString(rawRewardType)
		if err != nil {
			return nil, fmt.Errorf("convert reward type: %w", err)
		}

		costs[rewardType] = gloink.AmountFromInt(rawAmount)
	}

	return costs, nil
}
