package box

import (
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
)

type OpenBoxButtonPayload struct {
	ID ID
}

func OpenBoxButton(boxID ID) (button.Button, error) {
	btn, err := button.CreateButton(
		"Open",
		button.OperationOpenBox,
		false,
		button.StyleDefault,
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
		button.StyleDefault,
	)
}

func ShopButton(caption string) button.Button {
	return button.CreateButtonWithoutPayload(
		caption,
		button.OperationShop,
		false,
		button.StyleDefault,
	)
}

func CostButton(caption string) button.Button {
	return button.CreateButtonWithoutPayload(
		caption,
		button.OperationShop,
		true,
		button.StyleDefault,
	)
}

func InProgressButton(caption string) button.Button {
	return button.CreateButtonWithoutPayload(
		caption,
		button.OperationShop,
		true,
		button.StyleDanger,
	)
}

type ReadyToOpenButtonPayload struct {
	ID ID
}

func ReadyToOpenButton(caption string, boxID ID) (button.Button, error) {
	btn, err := button.CreateButton(
		caption,
		button.OperationBoxReadyToOpen,
		true,
		button.StyleSuccess,
		ReadyToOpenButtonPayload{
			ID: boxID,
		},
	)

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}
