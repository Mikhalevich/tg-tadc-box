package openbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

func (o *OpenBox) Open(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := o.transactor.Transaction(ctx, func(ctx context.Context) error {
		now := o.timeProvider.Now()

		_, isNewPlayer, err := o.getOrCreateUserProfile(ctx, chatID, now)
		if err != nil {
			return fmt.Errorf("get user profile: %w", err)
		}

		isAvailable, err := o.isNewBoxScheduleAvailable(ctx, chatID, now)
		if err != nil {
			return fmt.Errorf("is box schedule available: %w", err)
		}

		if !isAvailable {
			return nil
		}

		if err := o.scheduleBox(ctx, chatID, now, isNewPlayer); err != nil {
			return fmt.Errorf("schedule box: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

// getOrCreateUserProfile get or create user profile
// returns user, is user created flag and error.
func (o *OpenBox) getOrCreateUserProfile(
	ctx context.Context,
	chatID msginfo.ChatID,
	now time.Time,
) (player.Player, bool, error) {
	plr, err := o.repo.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return player.Player{}, false, fmt.Errorf("get player by chat id: %w", err)
		}

		plr, err := o.createPlayer(ctx, chatID, now)
		if err != nil {
			return player.Player{}, false, fmt.Errorf("create player: %w", err)
		}

		return plr, true, nil
	}

	return plr, false, nil
}

func (o *OpenBox) createPlayer(
	ctx context.Context,
	chatID msginfo.ChatID,
	createdAt time.Time,
) (player.Player, error) {
	plr := player.Player{
		ChatID:           chatID,
		CreatedAt:        createdAt,
		ProfileUpdatedAt: createdAt,
	}

	id, err := o.repo.InsertPlayer(ctx, plr)
	if err != nil {
		return player.Player{}, fmt.Errorf("insert player: %w", err)
	}

	plr.ID = player.IDFromInt(id)

	return plr, nil
}

func (o *OpenBox) isNewBoxScheduleAvailable(
	ctx context.Context,
	chatID msginfo.ChatID,
	now time.Time,
) (bool, error) {
	boxes, err := o.repo.GetBoxesByStatus(ctx, chatID, box.StatusInProgress)
	if err != nil {
		return false, fmt.Errorf("get active boxes: %w", err)
	}

	if len(boxes) > 0 {
		if err := o.notifier.ShowBoxInfo(ctx, boxes[0], boxes[0].AvailableAfter(now)); err != nil {
			return false, fmt.Errorf("show existing box info: %w", err)
		}

		return false, nil
	}

	return true, nil
}

func (o *OpenBox) scheduleBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	createdAt time.Time,
	isImmediate bool,
) error {
	newBox := createNormalBox(
		chatID,
		createdAt,
		calculateWaitPeriod(o.commonWaitPeriod, isImmediate),
	)

	boxID, err := o.repo.InsertBox(ctx, newBox)
	if err != nil {
		return fmt.Errorf("insert box: %w", err)
	}

	newBox.ID = box.IDFromInt(boxID)

	if err := o.notifier.ShowBoxInfo(ctx, newBox, newBox.AvailableAfter(createdAt)); err != nil {
		return fmt.Errorf("show new box info: %w", err)
	}

	return nil
}

func calculateWaitPeriod(standartPeriod time.Duration, isImmediate bool) time.Duration {
	if isImmediate {
		return 0
	}

	return standartPeriod
}

func createNormalBox(
	chatID msginfo.ChatID,
	createdAt time.Time,
	commonWaitPeriod time.Duration,
) box.Box {
	var (
		availableAt         = createdAt
		readyNotificationAt time.Time
	)

	if commonWaitPeriod == 0 {
		readyNotificationAt = createdAt
	} else {
		availableAt = availableAt.Add(commonWaitPeriod)
	}

	return box.Box{
		ChatID:              chatID,
		Status:              box.StatusInProgress,
		Type:                box.TypeNormal,
		CreatedAt:           createdAt,
		ReadyNotificationAt: readyNotificationAt,
		AvailableAt:         availableAt,
	}
}
