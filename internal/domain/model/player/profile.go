package player

import (
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Profile struct {
	OpenedBoxes OpenedBoxes
	Cards       CardsCollected
	Wallet      Wallet
}

// AbstractDuplicatesAll remove all duplicates from cards
// and add abstracted gloinks amount to wallet
// returns abstracted gloinks amount.
func (p *Profile) AbstractDuplicatesAll(
	abstractionCosts map[reward.RewardType]gloink.Amount,
) gloink.Amount {
	gloinksAmount := p.Cards.AbstractDuplicatesAll(abstractionCosts)
	p.Wallet.GloinksAmount += gloinksAmount

	return gloinksAmount
}

func (p *Profile) IsNoOpenCommonBoxes() bool {
	return p.OpenedBoxes.CommonCount() == 0
}
