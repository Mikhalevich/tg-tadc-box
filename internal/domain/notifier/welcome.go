package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (n *Notifier) Welcome(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	msg := `Welcome to the Amazing Digital Circus card collection bot.
Click button to receive your first reward.
`
	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypePlain,
			Text:   msg,
			Buttons: []button.ButtonRow{
				button.Row(box.GetCommonBoxButton("Free Box")),
			},
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
