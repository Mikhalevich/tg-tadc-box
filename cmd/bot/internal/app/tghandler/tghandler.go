package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type ButtonProvider interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
}

type cbHandler func(ctx context.Context, msg tgbot.BotMessage, btn *button.Button) error

type BoxScheduler interface {
	Schedule(ctx context.Context, chatID msginfo.ChatID) error
}

type ErrorNotifier interface {
	ParseError(ctx context.Context, chatID msginfo.ChatID, err error) error
}

type TGHandler struct {
	cbHanlers      map[button.Operation]cbHandler
	buttonProvider ButtonProvider
	boxScheduler   BoxScheduler
	errorNotifier  ErrorNotifier
}

func New(
	buttonProvider ButtonProvider,
	boxScheduler BoxScheduler,
	errorNotifier ErrorNotifier,
) *TGHandler {
	tgh := &TGHandler{
		buttonProvider: buttonProvider,
		boxScheduler:   boxScheduler,
		errorNotifier:  errorNotifier,
	}

	tgh.registerCBHandlers()

	return tgh
}

func (t *TGHandler) registerCBHandlers() {
	t.cbHanlers = map[button.Operation]cbHandler{}
}
