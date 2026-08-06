package player

import (
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

const (
	freeBoxAfterOpenedBoxCount = 10
)

type OpenedBoxes struct {
	Common    int
	Rare      int
	Epic      int
	Legendary int
}

func (ob *OpenedBoxes) Count() int {
	return ob.Common + ob.Rare + ob.Epic + ob.Legendary
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

func (ob *OpenedBoxes) IsFreeBoxAvailable(openedBoxType box.Type) (box.Type, bool) {
	count := ob.openedBoxCount(openedBoxType)

	if (count % freeBoxAfterOpenedBoxCount) != 0 {
		return box.TypeCommon, false
	}

	return box.TypeCommon, false
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
