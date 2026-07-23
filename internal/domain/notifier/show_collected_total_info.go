package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (n *Notifier) ShowCollectedTotalInfo(
	ctx context.Context,
	chatID msginfo.ChatID,
	infos []card.CollectedCardInfo,
) error {
	buttons, err := makeCollectedTotalInfoButtons(infos)
	if err != nil {
		return fmt.Errorf("make collected total info buttons: %w", err)
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:  chatID,
			Type:    msginfo.MessageTypeMarkdown,
			Text:    makeCollectedTotalInfoMsg(infos),
			Buttons: buttons,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func makeCollectedTotalInfoMsg(infos []card.CollectedCardInfo) string {
	lines := make([]string, 0, len(infos))

	for _, info := range infos {
		line := fmt.Sprintf("%s %d/%d", info.Type.String(), info.Collected, info.Total)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func makeCollectedTotalInfoButtons(infos []card.CollectedCardInfo) ([]button.ButtonRow, error) {
	buttons := make([]button.ButtonRow, 0, len(infos))

	for _, info := range infos {
		if info.Collected == 0 {
			continue
		}

		btn, err := card.PageButton(
			fmt.Sprintf("Show %s(%d)", info.Type.String(), info.Collected),
			info.Type,
			1,
		)
		if err != nil {
			return nil, fmt.Errorf("crate page button: %w", err)
		}

		buttons = append(buttons, button.Row(btn))
	}

	return buttons, nil
}
