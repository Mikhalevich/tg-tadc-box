package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
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

func (t *TGHandler) cbOpenBox(ctx context.Context, msg tgbot.BotMessage, btn *button.Button) error {
	payload, err := button.GetPayload[box.OpenBoxButtonPayload](*btn)
	if err != nil {
		return fmt.Errorf("get open box payload: %w", err)
	}

	if err := t.boxProcessor.OpenByID(ctx, msginfo.ChatIDFromInt64(msg.ChatID), payload.ID); err != nil {
		return fmt.Errorf("open box by id: %w", err)
	}

	return nil
}
