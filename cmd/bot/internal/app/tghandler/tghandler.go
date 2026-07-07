package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
)

type ButtonProvider interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
}

type cbHandler func(ctx context.Context, msg tgbot.BotMessage, btn *button.Button) error

type TGHandler struct {
	cbHanlers      map[button.Operation]cbHandler
	buttonProvider ButtonProvider
}

func New(
	buttonProvider ButtonProvider,
) *TGHandler {
	tgh := &TGHandler{
		buttonProvider: buttonProvider,
	}

	tgh.registerCBHandlers()

	return tgh
}

func (t *TGHandler) registerCBHandlers() {
	t.cbHanlers = map[button.Operation]cbHandler{}
}
