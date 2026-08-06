package box

import (
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
)

type OpenBoxButtonPayload struct {
	ID int
}

func OpenBoxButton(boxID int) (button.Button, error) {
	btn, err := button.CreateButton(
		"Open",
		button.OperationOpenBox,
		false,
		OpenBoxButtonPayload{
			ID: boxID,
		},
	)

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}

func GetCommonBoxButton(caption string) button.Button {
	return button.CreateButtonWithoutPayload(
		caption,
		button.OperationGetCommonBox,
		true,
	)
}

func ShopButton(caption string) button.Button {
	return button.CreateButtonWithoutPayload(
		caption,
		button.OperationShop,
		false,
	)
}
