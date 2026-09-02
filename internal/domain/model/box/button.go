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
		button.WithPayload(
			OpenBoxButtonPayload{
				ID: boxID,
			},
		),
	)

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}

func GetCommonBoxButton(caption string) button.Button {
	return button.MustCreateButton(
		caption,
		button.OperationGetCommonBox,
		button.WithDeleteAfterProcess(),
	)
}

func ShopButton(caption string) button.Button {
	return button.MustCreateButton(
		caption,
		button.OperationShop,
	)
}

func CostButton(caption string) button.Button {
	return button.MustCreateButton(
		caption,
		button.OperationShop,
		button.WithDeleteAfterProcess(),
	)
}

func InProgressButton(caption string) button.Button {
	return button.MustCreateButton(
		caption,
		button.OperationShop,
		button.WithDeleteAfterProcess(),
		button.WithStyle(button.StyleDanger),
	)
}

type ReadyToOpenButtonPayload struct {
	ID ID
}

func ReadyToOpenButton(caption string, boxID ID) (button.Button, error) {
	btn, err := button.CreateButton(
		caption,
		button.OperationBoxReadyToOpen,
		button.WithDeleteAfterProcess(),
		button.WithStyle(button.StyleSuccess),
		button.WithPayload(
			ReadyToOpenButtonPayload{
				ID: boxID,
			},
		),
	)

	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}
