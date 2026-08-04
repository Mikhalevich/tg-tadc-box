package player

type Profile struct {
	OpenedBoxes OpenedBoxes
	Cards       CardsCollected
	Wallet      Wallet
}

// AbstractDuplicatesAll remove all duplicates from cards
// and add abstracted gloinks amount to wallet
// returns abstracted gloinks amount.
func (p *Profile) AbstractDuplicatesAll() int {
	gloinksAmount := p.Cards.AbstractDuplicatesAll()
	p.Wallet.GloinksAmount += gloinksAmount

	return gloinksAmount
}
