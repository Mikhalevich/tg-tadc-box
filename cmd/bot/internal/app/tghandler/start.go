package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tgbot"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tghandler/internal/cmdargs"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

func (t *TGHandler) Start(
	ctx context.Context,
	msg tgbot.BotMessage,
	sender tgbot.MessageSender,
) error {
	operation, args := cmdargs.Parse(msg.Args)

	switch operation {
	case cmdargs.OperationJoin:
		if err := t.joinByLink.Join(
			ctx,
			msginfo.ChatIDFromInt64(msg.ChatID),
			referral.CodeFromString(args),
		); err != nil {
			return fmt.Errorf("join by link: %w", err)
		}

	case cmdargs.OperationNope:
		if err := t.notifier.Welcome(ctx, msginfo.ChatIDFromInt64(msg.ChatID)); err != nil {
			return fmt.Errorf("welcome: %w", err)
		}
	}

	return nil
}
