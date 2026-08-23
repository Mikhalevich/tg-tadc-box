package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Service) NoCollectedCards(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypePlain,
			Text:   "No cards collected",
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
