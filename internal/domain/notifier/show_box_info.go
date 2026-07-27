package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (n *Notifier) ShowBoxInfo(
	ctx context.Context,
	domBox box.Box,
	availableAfter time.Duration,
) error {
	if availableAfter > 0 {
		if err := n.sendBoxIsNotAvailableYet(ctx, domBox, availableAfter); err != nil {
			return fmt.Errorf("send box is not available yet: %w", err)
		}

		return nil
	}

	if err := n.sendBoxIsAvailable(ctx, domBox); err != nil {
		return fmt.Errorf("send box is available: %w", err)
	}

	return nil
}

func (n *Notifier) sendBoxIsNotAvailableYet(
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

func (n *Notifier) sendBoxIsAvailable(
	ctx context.Context,
	domBox box.Box,
) error {
	openBoxBtn, err := box.OpenBoxButton(domBox.ID.Int())
	if err != nil {
		return fmt.Errorf("open box button: %w", err)
	}

	payload, err := n.imageProvider.Chest(ctx)
	if err != nil {
		return fmt.Errorf("receive chest paylod: %w", err)
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:  domBox.ChatID,
			Type:    msginfo.MessageTypePNG,
			Payload: payload,
			Buttons: []button.ButtonRow{
				{
					openBoxBtn,
				},
			},
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
