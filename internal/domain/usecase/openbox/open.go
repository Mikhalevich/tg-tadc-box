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

		if _, err := o.processScheduleBox(
			ctx,
			chatID,
			box.TypeCommon,
			isNoOpenBoxes(profile),
		); err != nil {
			return fmt.Errorf("process schedule box: %w", err)
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
