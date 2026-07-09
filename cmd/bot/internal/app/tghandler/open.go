package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (t *TGHandler) Open(
	ctx context.Context,
	msg tgbot.BotMessage,
	sender tgbot.MessageSender,
) error {
	if err := t.boxScheduler.Schedule(ctx, msginfo.ChatIDFromInt64(msg.ChatID)); err != nil {
		return fmt.Errorf("schedule box: %w", err)
	}

	return nil
}
