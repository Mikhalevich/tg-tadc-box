package gloink

import (
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
)

type BuyBoxButtonPayload struct {
	BoxType box.Type
}

func BuyBoxButton(
	caption string,
	boxType box.Type,
) (button.Button, error) {
	btn, err := button.CreateButton(
		caption,
		button.OperationBuyBox,
		true,
		BuyBoxButtonPayload{
			BoxType: boxType,
		},
	)
	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}
