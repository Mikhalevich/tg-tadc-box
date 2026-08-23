package showreadybox

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type BoxService interface {
	GetBoxByID(ctx context.Context, id box.ID) (box.Box, error)
}

type NotificationService interface {
	ShowReadyToOpenBox(ctx context.Context, domBox box.Box) error
}

type TimeProvider interface {
	Now() time.Time
}

type ShowReadyToOpenBox struct {
	boxService          BoxService
	notificationService NotificationService
	timeProvider        TimeProvider
}

func New(
	boxService BoxService,
	notificationService NotificationService,
	timeProvider TimeProvider,
) *ShowReadyToOpenBox {
	return &ShowReadyToOpenBox{
		boxService:          boxService,
		notificationService: notificationService,
		timeProvider:        timeProvider,
	}
}

func (s *ShowReadyToOpenBox) ShowReadyToOpenBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxID box.ID,
) error {
	readyBox, err := s.boxService.GetBoxByID(ctx, boxID)
	if err != nil {
		return fmt.Errorf("get box by id: %w", err)
	}

	if err := readyBox.IsReadyToOpen(s.timeProvider.Now()); err != nil {
		return fmt.Errorf("check box is ready for open: %w", err)
	}

	if err := s.notificationService.ShowReadyToOpenBox(ctx, readyBox); err != nil {
		return fmt.Errorf("show ready to open notification: %w", err)
	}

	return nil
}
