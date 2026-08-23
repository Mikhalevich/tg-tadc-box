package notificationsvc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Service) ShowCardPageTotal(
	ctx context.Context,
	chatID msginfo.ChatID,
	infos []card.CardPageTotal,
	abstractDuplicatesAmount gloink.Amount,
) error {
	buttons, err := makeCardPageTotalButtons(infos, abstractDuplicatesAmount)
	if err != nil {
		return fmt.Errorf("make collected total info buttons: %w", err)
	}

	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:  chatID,
			Type:    msginfo.MessageTypeMarkdown,
			Text:    makeCardPageTotalMsg(infos),
			Buttons: buttons,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func makeCardPageTotalMsg(infos []card.CardPageTotal) string {
	lines := make([]string, 0, len(infos))

	for _, info := range infos {
		line := fmt.Sprintf("%s *%d*/%d", info.Type.String(), info.Collected, info.Total)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func calculateButtonsCount(
	infos []card.CardPageTotal,
	abstractDuplicatesAmount gloink.Amount,
) int {
	if abstractDuplicatesAmount > 0 {
		return len(infos) + 1
	}

	return len(infos)
}

func makeCardPageTotalButtons(
	infos []card.CardPageTotal,
	abstractDuplicatesAmount gloink.Amount,
) ([]button.ButtonRow, error) {
	buttons := make([]button.ButtonRow, 0,
		calculateButtonsCount(infos, abstractDuplicatesAmount))

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

	if abstractDuplicatesAmount > 0 {
		buttons = append(buttons,
			button.Row(
				card.AbstractDuplicatesAllButton(
					fmt.Sprintf("Abstract duplicates for %d gloinks",
						abstractDuplicatesAmount),
				),
			),
		)
	}

	return buttons, nil
}
