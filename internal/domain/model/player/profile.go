package player

type Profile struct {
	OpenedBoxes OpenedBoxes
	Cards       CardsCollected
}

type OpenedBoxes struct {
	Common int
}

func (ob OpenedBoxes) Count() int {
	return ob.Common
}
