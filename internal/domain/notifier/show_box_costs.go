package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (n *Notifier) ShowBoxCosts(
	ctx context.Context,
	chatID msginfo.ChatID,
	wallet player.Wallet,
	costs []gloink.BoxCost,
	inProgressBoxes map[box.Type]box.Box,
) error {
	var buttonRows []button.ButtonRow
	for _, cost := range costs {
		_, isBoxAlreadyInProgress := inProgressBoxes[cost.Type]

		btn, err := gloink.BuyBoxButton(
			messageBoxAlreadyInProgress(
				messageForBoxAmount(cost.Type, cost.Amount),
				isBoxAlreadyInProgress,
			),
			cost.Type,
		)
		if err != nil {
			return fmt.Errorf("create buy box button: %w", err)
		}

		buttonRows = append(buttonRows, button.Row(btn))
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:  chatID,
			Type:    msginfo.MessageTypeMarkdown,
			Text:    fmt.Sprintf("Gloinks available *%d*", wallet.GloinksAmount.Int()),
			Buttons: buttonRows,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func messageForBoxAmount(
	boxType box.Type,
	amount gloink.Amount,
) string {
	if amount.Int() == 0 {
		return "Free"
	}

	return fmt.Sprintf("%s %d gloinks", boxType.String(), amount.Int())
}

func messageBoxAlreadyInProgress(
	msg string,
	isInProgress bool,
) string {
	if !isInProgress {
		return msg
	}

	return fmt.Sprintf("%s (⌛️)", msg)
}
