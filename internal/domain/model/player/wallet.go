package player

import (
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

type Wallet struct {
	GloinksAmount gloink.Amount
}

func (w *Wallet) DecreaseGloinks(amount gloink.Amount) error {
	if amount > w.GloinksAmount {
		return perror.InvalidParam("not enough gloinks")
	}

	w.GloinksAmount -= amount

	return nil
}

func (w *Wallet) IncreaseGloinks(amount gloink.Amount) {
	w.GloinksAmount += amount
}
