package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
)

func (t *TGHandler) DefaultHandler(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	return nil
}
