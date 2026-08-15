package notifier

import (
	"context"
	"fmt"
	"time"

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
	inProgressBoxes map[box.Type]box.InProgressBox,
) error {
	buttons, err := n.makeBoxCostsButtons(costs, inProgressBoxes)
	if err != nil {
		return fmt.Errorf("make box costs buttons: %w", err)
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:  chatID,
			Type:    msginfo.MessageTypeMarkdown,
			Text:    fmt.Sprintf("Gloinks available *%d*", wallet.GloinksAmount.Int()),
			Buttons: buttons,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (n *Notifier) makeBoxCostsButtons(
	costs []gloink.BoxCost,
	inProgressBoxes map[box.Type]box.InProgressBox,
) ([]button.ButtonRow, error) {
	var buttonRows []button.ButtonRow
	for _, cost := range costs {
		inProgressBox, isBoxAlreadyInProgress := inProgressBoxes[cost.Type]

		btn, err := n.createBoxCostButton(cost, isBoxAlreadyInProgress, inProgressBox.AvailableAfter)
		if err != nil {
			return nil, fmt.Errorf("create button: %w", err)
		}

		buttonRows = append(buttonRows, button.Row(btn))
	}

	return buttonRows, nil
}

func (n *Notifier) createBoxCostButton(
	cost gloink.BoxCost,
	isInProgress bool,
	availableAfter time.Duration,
) (button.Button, error) {
	if isInProgress {
		return button.CreateNoOperationButton(
			n.messageBoxAlreadyInProgress(cost.Type, availableAfter),
		), nil
	}

	buyButton, err := gloink.BuyBoxButton(
		messageForBoxAmount(cost.Type, cost.Amount),
		cost.Type,
	)

	if err != nil {
		return button.Button{}, fmt.Errorf("create buy box button: %w", err)
	}

	return buyButton, nil
}

func messageForBoxAmount(
	boxType box.Type,
	amount gloink.Amount,
) string {
	if amount.Int() == 0 {
		return "Free"
	}

	return fmt.Sprintf("%s %d gloinks", boxType.Pretty(), amount.Int())
}

func (n *Notifier) messageBoxAlreadyInProgress(
	boxType box.Type,
	availableAfter time.Duration,
) string {
	if availableAfter > 0 {
		return fmt.Sprintf("%s ⌛️ %s", boxType.Pretty(), n.parseDuration(availableAfter))
	}

	return fmt.Sprintf("%s is ready to open", boxType.Pretty())
}
