package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (n *Notifier) NoCollectedCards(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypePlain,
			Text:   "No cards collected",
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
