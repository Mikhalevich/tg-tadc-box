package card

import (
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type PageButtonPayload struct {
	Type reward.RewardType
	Page int
}

func PageButton(caption string, rewardType reward.RewardType, page int) (button.Button, error) {
	btn, err := button.CreateButton(
		caption,
		button.OperationCardPage,
		false,
		button.StyleDefault,
		PageButtonPayload{
			Type: rewardType,
			Page: page,
		},
	)
	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}

func TotalButton(caption string) button.Button {
	return button.CreateButtonWithoutPayload(
		caption,
		button.OperationCardTotal,
		true,
		button.StyleDefault,
	)
}

func AbstractDuplicatesAllButton(caption string) button.Button {
	return button.CreateButtonWithoutPayload(
		caption,
		button.OperationAbstractDuplicatesAll,
		false,
		button.StyleDefault,
	)
}

type AbstractCardPayload struct {
	Type  reward.RewardType
	Pos   int
	Count int
}

func AbstractCardButton(
	caption string,
	rewardType reward.RewardType,
	pos int,
	count int,
) (button.Button, error) {
	btn, err := button.CreateButton(
		caption,
		button.OperationAbstractCard,
		false,
		button.StyleDefault,
		AbstractCardPayload{
			Type:  rewardType,
			Pos:   pos,
			Count: count,
		},
	)

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}
