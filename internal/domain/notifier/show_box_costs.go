package notifier

import (
	"context"
	"fmt"

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
) error {
	var buttonRows []button.ButtonRow
	for _, cost := range costs {
		btn, err := gloink.BuyBoxButton(
			fmt.Sprintf("Buy %s for %d gloinks", cost.Type.String(), cost.Amount.Int()),
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
