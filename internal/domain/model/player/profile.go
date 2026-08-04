package player

type Profile struct {
	OpenedBoxes OpenedBoxs
	Cards       CardsCollected
}

type OpenedBoxs struct {
	Common int
}
