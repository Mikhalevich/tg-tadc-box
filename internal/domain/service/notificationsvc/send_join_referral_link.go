package notificationsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

func (s *Service) SendJoinReferralLink(
	ctx context.Context,
	chatID msginfo.ChatID,
	code referral.Code,
) error {
	if err := s.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Type:   msginfo.MessageTypePlain,
			Text:   makeJoinReferralLink(code.String()),
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func makeJoinReferralLink(code string) string {
	return fmt.Sprintf("https://t.me/tadcBoxBot?start=join_%s", code)
}
