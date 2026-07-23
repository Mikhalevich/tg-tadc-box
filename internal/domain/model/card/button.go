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
