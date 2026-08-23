package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

//nolint:funlen
func (s *Service) ShowCardPage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	rew reward.Reward,
	count int,
	page int,
	maxPage int,
	firstPage card.CardPage,
	previousPage card.CardPage,
	nextPage card.CardPage,
	lastPage card.CardPage,
) error {
	payload, err := s.imageProvider.Reward(ctx, rew)
	if err != nil {
		return fmt.Errorf("receive image payload: %w", err)
	}

	navigationButtons, err := makeCardPageNavigationButtons(
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

	cmdButtons, err := makeCardPageCommandButtons(rew.Type, page, count)
	if err != nil {
		return fmt.Errorf("make card page command buttons: %w", err)
	}

	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:     chatID,
			ReplyMsgID: messageID,
			Type:       msginfo.MessageTypePNG,
			Payload:    payload,
			Text:       fmt.Sprintf("%d/%d x%d", page, maxPage, count),
			Buttons: []button.ButtonRow{
				cmdButtons,
				navigationButtons,
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

func makeCardPageCommandButtons(
	rewardType reward.RewardType,
	page int,
	count int,
) (button.ButtonRow, error) {
	backBtn := card.TotalButton("Back")

	if count <= 1 {
		return button.Row(backBtn), nil
	}

	abstractOneBtn, err := card.AbstractCardButton(
		"Abstract 1",
		rewardType, page, 1)
	if err != nil {
		return nil, fmt.Errorf("abstract one button: %w", err)
	}

	//nolint:mnd
	if count == 2 {
		return button.Row(backBtn, abstractOneBtn), nil
	}

	abstractAllBtn, err := card.AbstractCardButton(
		fmt.Sprintf("Abstract %d", count-1),
		rewardType, page, count-1)
	if err != nil {
		return nil, fmt.Errorf("abstract all: %w", err)
	}

	return button.Row(backBtn, abstractOneBtn, abstractAllBtn), nil
}

func makeCardPageNavigationButtons(
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
