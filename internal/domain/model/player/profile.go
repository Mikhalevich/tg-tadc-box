package player

import (
	"fmt"

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

func (p *Profile) AbstractByPos(
	rewardType reward.RewardType,
	pos int,
	count int,
	abstractionCosts map[reward.RewardType]gloink.Amount,
) (gloink.Amount, error) {
	gloinksAmount, err := p.Cards.AbstractByPos(
		rewardType,
		pos,
		count,
		abstractionCosts,
	)
	if err != nil {
		return 0, fmt.Errorf("abstract by pos: %w", err)
	}

	p.Wallet.GloinksAmount += gloinksAmount

	return gloinksAmount, nil
}

func (p *Profile) IsNoOpenCommonBoxes() bool {
	return p.OpenedBoxes.CommonCount() == 0
}
