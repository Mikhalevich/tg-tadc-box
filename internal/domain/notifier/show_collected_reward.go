package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (n *Notifier) ShowCollectedReward(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	rew reward.Reward,
	count int,
	previousPage card.CardPage,
	nextPage card.CardPage,
) error {
	payload, err := n.imageProvider.Reward(ctx, rew)
	if err != nil {
		return fmt.Errorf("receive image payload: %w", err)
	}

	buttons, err := makeCardPageButtons(rew.Type, previousPage, nextPage)
	if err != nil {
		return fmt.Errorf("make cards buttons: %w", err)
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:     chatID,
			ReplyMsgID: messageID,
			Type:       msginfo.MessageTypeEditPNG,
			Payload:    payload,
			Text:       fmt.Sprintf("x%d", count),
			Buttons: []button.ButtonRow{
				button.Row(card.TotalButton("Back")),
				buttons,
			},
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func makeCardPageButtons(
	rewardType reward.RewardType,
	previousPage card.CardPage,
	nextPage card.CardPage,
) (button.ButtonRow, error) {
	var buttons button.ButtonRow

	if previousPage.IsValid {
		prevBtn, err := card.PageButton("<", rewardType, previousPage.Page)
		if err != nil {
			return nil, fmt.Errorf("prev button: %w", err)
		}

		buttons = append(buttons, prevBtn)
	}

	if nextPage.IsValid {
		nextBtn, err := card.PageButton(">", rewardType, nextPage.Page)
		if err != nil {
			return nil, fmt.Errorf("next button: %w", err)
		}

		buttons = append(buttons, nextBtn)
	}

	return buttons, nil
}
