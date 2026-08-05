package gloink

import (
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type BuyBoxButtonPayload struct {
	Type reward.RewardType
}

func BuyBoxButton(
	caption string,
	rewardType reward.RewardType,
) (button.Button, error) {
	btn, err := button.CreateButton(
		caption,
		button.OperationBuyBox,
		true,
		BuyBoxButtonPayload{
			Type: rewardType,
		},
	)
	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}
