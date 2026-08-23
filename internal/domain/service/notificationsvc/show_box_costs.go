package notificationsvc

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

func (s *Service) ShowBoxCosts(
	ctx context.Context,
	chatID msginfo.ChatID,
	wallet player.Wallet,
	costs []gloink.BoxCost,
	inProgressBoxes map[box.Type]box.InProgressBox,
) error {
	buttons, err := s.makeBoxCostsButtons(costs, inProgressBoxes)
	if err != nil {
		return fmt.Errorf("make box costs buttons: %w", err)
	}

	if err := s.sender.SendMessage(
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

func (s *Service) makeBoxCostsButtons(
	costs []gloink.BoxCost,
	inProgressBoxes map[box.Type]box.InProgressBox,
) ([]button.ButtonRow, error) {
	var buttonRows []button.ButtonRow
	for _, cost := range costs {
		inProgressBox, isBoxAlreadyInProgress := inProgressBoxes[cost.Type]

		btn, err := s.createBoxCostButton(cost, isBoxAlreadyInProgress, inProgressBox)
		if err != nil {
			return nil, fmt.Errorf("create button: %w", err)
		}

		buttonRows = append(buttonRows, button.Row(btn))
	}

	return buttonRows, nil
}

func (s *Service) createBoxCostButton(
	cost gloink.BoxCost,
	isInProgress bool,
	inProgressBox box.InProgressBox,
) (button.Button, error) {
	if isInProgress {
		if inProgressBox.AvailableAfter > 0 {
			return box.InProgressButton(
				s.msgBoxInProgress(cost.Type, inProgressBox.AvailableAfter),
			), nil
		}

		readyToOpenBtn, err := box.ReadyToOpenButton(
			msgBoxInProgress(cost.Type),
			inProgressBox.Box.ID,
		)

		if err != nil {
			return button.Button{}, fmt.Errorf("create ready to open button: %w", err)
		}

		return readyToOpenBtn, nil
	}

	buyButton, err := gloink.BuyBoxButton(
		msgBoxCost(cost.Type, cost.Amount),
		cost.Type,
	)

	if err != nil {
		return button.Button{}, fmt.Errorf("create buy box button: %w", err)
	}

	return buyButton, nil
}

func msgBoxCost(
	boxType box.Type,
	amount gloink.Amount,
) string {
	if amount.Int() == 0 {
		return "Free"
	}

	return fmt.Sprintf("%s %d gloinks", boxType.Pretty(), amount.Int())
}

func msgBoxInProgress(boxType box.Type) string {
	return fmt.Sprintf("%s is ready", boxType.Pretty())
}

func (s *Service) msgBoxInProgress(
	boxType box.Type,
	availableAfter time.Duration,
) string {
	return fmt.Sprintf("%s ⌛️ %s", boxType.Pretty(), s.parseDuration(availableAfter))
}
