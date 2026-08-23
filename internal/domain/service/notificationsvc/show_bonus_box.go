package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (s *Service) ShowBonusBox(
	ctx context.Context,
	bonusBox box.Box,
) error {
	msg := fmt.Sprintf(
		"Received *%s* bonus box for opening x%d *%s* boxes",
		bonusBox.Type.String(),
		bonusBox.Meta.BonusBox.Attempts,
		bonusBox.Meta.BonusBox.Type.String(),
	)
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: bonusBox.ChatID,
			Type:   msginfo.MessageTypeMarkdown,
			Text:   msg,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
