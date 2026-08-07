package player

import (
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

type OpenedBoxes struct {
	Common    int
	Rare      int
	Epic      int
	Legendary int
}

func (ob *OpenedBoxes) CommonCount() int {
	return ob.Common
}

func (ob *OpenedBoxes) Add(boxType box.Type) {
	switch boxType {
	case box.TypeCommon:
		ob.Common++

	case box.TypeRare:
		ob.Rare++

	case box.TypeEpic:
		ob.Epic++

	case box.TypeLegendary:
		ob.Legendary++
	}
}

// IsBonusBoxAvailable check for free box availability
// returns free box type and availability flag.
func (ob *OpenedBoxes) IsBonusBoxAvailable(
	openedBoxType box.Type,
	bonusBoxAfterOpenedBoxCount int,
) (box.Type, bool) {
	count := ob.openedBoxCount(openedBoxType)

	if (count % bonusBoxAfterOpenedBoxCount) != 0 {
		return box.TypeCommon, false
	}

	bonusBoxType := openedBoxType.Next()

	return bonusBoxType, bonusBoxType.IsValid()
}

func (ob *OpenedBoxes) openedBoxCount(boxType box.Type) int {
	switch boxType {
	case box.TypeCommon:
		return ob.Common

	case box.TypeRare:
		return ob.Rare

	case box.TypeEpic:
		return ob.Epic

	case box.TypeLegendary:
		return ob.Legendary
	}

	return 0
}
