package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (n *Notifier) ShowGloinksWalletAfterAbstraction(
	ctx context.Context,
	chatID msginfo.ChatID,
	wallet player.Wallet,
	abstractedGloinksAmount int,
) error {
	var (
		msgTemplate = `Duplicates abstracted for *%d* gloinks
Your gloinks amount is *%d*
`
		msg = fmt.Sprintf(msgTemplate, abstractedGloinksAmount, wallet.GloinksAmount)
	)

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypeMarkdown,
			Text:   msg,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
