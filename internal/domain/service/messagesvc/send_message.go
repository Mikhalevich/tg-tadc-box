package messagesvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Service) SendMessage(
	ctx context.Context,
	msg msginfo.Message,
) error {
	inlineButtons, err := s.SetButtonRows(ctx, msg.Buttons...)
	if err != nil {
		return fmt.Errorf("set button rows: %w", err)
	}

	if err := s.sender.SendMessage(ctx, msginfo.SenderMessage{
		ChatID:     msg.ChatID,
		ReplyMsgID: msg.ReplyMsgID,
		Text:       msg.Text,
		Type:       msg.Type,
		Payload:    msg.Payload,
		Buttons:    inlineButtons,
	}); err != nil {
		return fmt.Errorf("sender send message: %w", err)
	}

	return nil
}
