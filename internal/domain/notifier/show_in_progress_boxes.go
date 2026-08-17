package notifier

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (n *Notifier) ShowInProgressBoxes(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxes []box.InProgressBox,
) error {
	inProgressBoxMsgLines := make([]string, 0, len(boxes))

	for _, inProgressBox := range boxes {
		if inProgressBox.AvailableAfter > 0 {
			msgLine := fmt.Sprintf("*%s* box will be available after *%s*",
				inProgressBox.Box.Type.String(),
				n.parseDuration(inProgressBox.AvailableAfter),
			)

			inProgressBoxMsgLines = append(inProgressBoxMsgLines, msgLine)

			continue
		}

		if err := n.sendBoxIsAvailable(ctx, inProgressBox.Box); err != nil {
			return fmt.Errorf("send box is available: %w", err)
		}
	}

	if len(inProgressBoxMsgLines) == 0 {
		return nil
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Text:   strings.Join(inProgressBoxMsgLines, "\n"),
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
	openBoxBtn, err := box.OpenBoxButton(domBox.ID)
	if err != nil {
		return fmt.Errorf("open box button: %w", err)
	}

	payload, err := n.imageProvider.Chest(ctx, domBox.Type)
	if err != nil {
		return fmt.Errorf("receive chest paylod: %w", err)
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:  domBox.ChatID,
			Type:    msginfo.MessageTypePNG,
			Text:    makeBonusBoxDescription(domBox.Meta),
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

func makeBonusBoxDescription(meta box.Meta) string {
	if !meta.BonusBox.IsValid {
		return ""
	}

	return fmt.Sprintf(
		"Bonus box for opening x%d *%s* boxes",
		meta.BonusBox.Attempts,
		meta.BonusBox.Type.String(),
	)
}

func (n *Notifier) parseDuration(d time.Duration) string {
	return n.escaper.EscapeMarkdown(d.Truncate(time.Second).String())
}
