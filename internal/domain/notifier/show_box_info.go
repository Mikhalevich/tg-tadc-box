package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (n *Notifier) ShowBoxInfo(
	ctx context.Context,
	domBox box.Box,
	availableAfter time.Duration,
) error {
	msg := fmt.Sprintf("Box will be available after *%s*",
		n.escaper.EscapeMarkdown(availableAfter.Truncate(time.Second).String()))

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: domBox.ChatID,
			Text:   msg,
			Type:   msginfo.MessageTypeMarkdown,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
