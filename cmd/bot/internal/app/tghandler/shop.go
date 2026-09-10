package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tgbot"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (t *TGHandler) Shop(
	ctx context.Context,
	msg tgbot.BotMessage,
	sender tgbot.MessageSender,
) error {
	if err := t.shop.ViewBoxes(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
	); err != nil {
		return fmt.Errorf("shop: %w", err)
	}

	return nil
}
