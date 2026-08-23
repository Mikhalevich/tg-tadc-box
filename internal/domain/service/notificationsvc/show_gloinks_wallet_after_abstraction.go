package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (s *Service) ShowGloinksWalletAfterAbstraction(
	ctx context.Context,
	chatID msginfo.ChatID,
	wallet player.Wallet,
	abstractedAmount gloink.Amount,
) error {
	var (
		msgTemplate = `Duplicates abstracted for *%d* gloinks
Your gloinks amount is *%d*
`
		msg = fmt.Sprintf(msgTemplate, abstractedAmount, wallet.GloinksAmount)
	)

	if err := s.sender.SendMessage(
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
