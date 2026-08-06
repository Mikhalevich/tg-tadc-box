package player

import (
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
)

type Profile struct {
	OpenedBoxes OpenedBoxes
	Cards       CardsCollected
	Wallet      Wallet
}

// AbstractDuplicatesAll remove all duplicates from cards
// and add abstracted gloinks amount to wallet
// returns abstracted gloinks amount.
func (p *Profile) AbstractDuplicatesAll() gloink.Amount {
	gloinksAmount := p.Cards.AbstractDuplicatesAll()
	p.Wallet.GloinksAmount += gloinksAmount

	return gloinksAmount
}

func (p *Profile) IsNoOpenCommonBoxes() bool {
	return p.OpenedBoxes.CommonCount() == 0
}
