package openbox

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (o *OpenBox) Open(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		profile, err := o.playerProvider.GetPlayerByChatID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("get user profile: %w", err)
		}

		inProgressBoxes, err := o.inProgressBoxesByType(ctx, chatID, box.TypeNormal)
		if err != nil {
			return fmt.Errorf("get in progress boxes by type %s: %w", box.TypeNormal, err)
		}

		now := o.timeProvider.Now()

		if len(inProgressBoxes) != 0 {
			if err := o.notifier.ShowBoxInfo(ctx, inProgressBoxes[0], inProgressBoxes[0].AvailableAfter(now)); err != nil {
				return fmt.Errorf("show existing box info: %w", err)
			}

			return nil
		}

		if err := o.scheduleBox(
			ctx,
			chatID,
			box.TypeNormal,
			now,
			isNoOpenBoxes(profile),
		); err != nil {
			return fmt.Errorf("schedule box: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func isNoOpenBoxes(profile player.Player) bool {
	return profile.Profile.OpenedBoxes.Count() == 0
}
