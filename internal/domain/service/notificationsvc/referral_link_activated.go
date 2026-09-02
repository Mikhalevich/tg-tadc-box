package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Service) ReferralLinkActivated(
	ctx context.Context,
	chatID msginfo.ChatID,
	gloinksReceived gloink.Amount,
) error {
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypeMarkdown,
			Text:   fmt.Sprintf("Share link activated, *%d* gloinks received", gloinksReceived.Int()),
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
