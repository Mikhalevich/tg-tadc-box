package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (t *TGHandler) DefaultCallbackQuery(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if msg.Data == "" {
		return nil
	}

	btn, err := t.buttonProvider.GetButton(ctx, button.IDFromString(msg.Data))
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			sender.SendMessage(ctx, msg.ChatID, "Button expired")

			return nil
		}

		return fmt.Errorf("get button: %w", err)
	}

	hndlr, ok := t.cbHanlers[btn.Operation]
	if !ok {
		return fmt.Errorf("invalid operation %v", btn.Operation)
	}

	if err := hndlr(ctx, msg, btn); err != nil {
		return fmt.Errorf("process cb handler: %w", err)
	}

	if btn.IsDeleteAfterProcess {
		sender.DeleteMessage(ctx, msg.ChatID, msg.MessageID)
	}

	return nil
}
