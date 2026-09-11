package app

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tgbot"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tghandler"
	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/config"
	"github.com/Mikhalevich/tg-tadc-box/internal/infra/logger"
	"github.com/Mikhalevich/tg-tadc-box/internal/infra/tracing"
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

	tbot, err := tgbot.New(
		botCfg.Token,
		tgbot.WithWebHookToken(botCfg.WebHookToken),
		tgbot.WithReadinessProbe(readinessProbe),
		tgbot.WithLivenessProbe(livenessProbe),
		tgbot.WithNewTracerFn(func() tgbot.Tracer {
			return newTracer(logger.FromContext(ctx))
		}),
	)
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	makeRoutes(tbot, botHandler)

	if err := tbot.Start(ctx); err != nil {
		return fmt.Errorf("bot start: %w", err)
	}

	return nil
}

type tracer struct {
	log     logger.Logger
	endSpan func()
}

func newTracer(log logger.Logger) *tracer {
	return &tracer{
		log: log,
	}
}

func (t *tracer) Before(
	ctx context.Context,
	pattern string,
	msg tgbot.BotMessage,
) context.Context {
	ctx, span := tracing.StartSpanName(ctx, pattern)

	t.endSpan = func() {
		span.End()
	}

	log := t.log.WithContext(ctx).
		WithField("endpoint", pattern).
		WithField("bot_message", msg)

	return logger.WithLogger(ctx, log)
}

func (t *tracer) OnSuccess(ctx context.Context) {
}

func (t *tracer) OnError(ctx context.Context, err error) {
	logger.FromContext(ctx).
		WithError(err).
		Error("error while processing message")
}

func (t *tracer) After(ctx context.Context) {
	t.endSpan()
}
