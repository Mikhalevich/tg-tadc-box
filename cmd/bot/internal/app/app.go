package app

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tghandler"
	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/config"
	"github.com/Mikhalevich/tg-tadc-box/internal/infra/logger"
)

func Start(
	ctx context.Context,
	botCfg config.Bot,
	buttonProvider tghandler.ButtonProvider,
	boxScheduler tghandler.BoxProcessor,
	cardViewer tghandler.CardViewer,
	cardAbstracter tghandler.CardAbstracter,
	shop tghandler.Shop,
	notifier tghandler.Notifier,
	errorNotifier tghandler.ErrorNotifier,
) error {
	var (
		botHandler = tghandler.New(
			buttonProvider,
			boxScheduler,
			cardViewer,
			cardAbstracter,
			shop,
			notifier,
			errorNotifier,
		)
	)

	tbot, err := tgbot.New(botCfg.Token, botCfg.WebHookToken, logger.FromContext(ctx))
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	makeRoutes(tbot, botHandler)

	if err := tbot.Start(ctx); err != nil {
		return fmt.Errorf("bot start: %w", err)
	}

	return nil
}
