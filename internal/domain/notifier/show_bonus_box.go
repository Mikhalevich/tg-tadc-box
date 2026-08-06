package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (n *Notifier) ShowBonusBox(
	ctx context.Context,
	bonusBox box.Box,
) error {
	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: bonusBox.ChatID,
			Type:   msginfo.MessageTypeMarkdown,
			Text:   fmt.Sprintf("Received *%s* bonus box", bonusBox.Type),
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
