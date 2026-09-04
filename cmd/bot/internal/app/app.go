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
	boxOpenByID tghandler.BoxOpenByID,
	boxShowReady tghandler.BoxShowReadyToOpen,
	cardViewer tghandler.CardViewer,
	cardAbstracter tghandler.CardAbstracter,
	shop tghandler.Shop,
	likeProcessor tghandler.LikeProcessor,
	inviteSender tghandler.SendInviteLink,
	joinByLink tghandler.JoinByLink,
	notifier tghandler.Notifier,
	errorNotifier tghandler.ErrorNotifier,
	readinessProbe tgbot.Probe,
	livenessProbe tgbot.Probe,
) error {
	var (
		botHandler = tghandler.New(
			buttonProvider,
			boxOpenByID,
			boxShowReady,
			cardViewer,
			cardAbstracter,
			shop,
			likeProcessor,
			inviteSender,
			joinByLink,
			notifier,
			errorNotifier,
		)
	)

	tbot, err := tgbot.New(botCfg.Token, botCfg.WebHookToken, logger.FromContext(ctx))
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	makeRoutes(tbot, botHandler)

	tbot.SetReadinessProbe(readinessProbe)
	tbot.SetLivenessProbe(livenessProbe)

	if err := tbot.Start(ctx); err != nil {
		return fmt.Errorf("bot start: %w", err)
	}

	return nil
}
