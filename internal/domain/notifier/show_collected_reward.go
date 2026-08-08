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
	page int,
	maxPage int,
	previousPage card.CardPage,
	nextPage card.CardPage,
	firstPage card.CardPage,
	lastPage card.CardPage,
) error {
	payload, err := n.imageProvider.Reward(ctx, rew)
	if err != nil {
		return fmt.Errorf("receive image payload: %w", err)
	}

	buttons, err := makeCardPageButtons(
		rew.Type,
		[]cardPageWithCaption{
			{
				Page:    firstPage,
				Caption: "<<",
			},
			{
				Page:    previousPage,
				Caption: "<",
			},
			{
				Page:    nextPage,
				Caption: ">",
			},
			{
				Page:    lastPage,
				Caption: ">>",
			},
		},
	)
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
			Text:       fmt.Sprintf("%d/%d x%d", page, maxPage, count),
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

type cardPageWithCaption struct {
	Page    card.CardPage
	Caption string
}

func makeCardPageButtons(
	rewardType reward.RewardType,
	pages []cardPageWithCaption,
) (button.ButtonRow, error) {
	var buttons button.ButtonRow

	for _, page := range pages {
		if !page.Page.IsValid {
			continue
		}

		btn, err := card.PageButton(page.Caption, rewardType, page.Page.Page)
		if err != nil {
			return nil, fmt.Errorf("page button: %w", err)
		}

		buttons = append(buttons, btn)
	}

	return buttons, nil
}
