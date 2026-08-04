package player

type OpenedBoxes struct {
	Common int
}

func (ob OpenedBoxes) Count() int {
	return ob.Common
}
