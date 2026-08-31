package invitelink

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

type ReferralService interface {
	CodeByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (referral.Code, error)
}

type Notifier interface {
	SendJoinReferralLink(
		ctx context.Context,
		chatID msginfo.ChatID,
		code referral.Code,
	) error
}

type InviteLink struct {
	referralService ReferralService
	notifier        Notifier
}

func New(
	referralService ReferralService,
	notifier Notifier,
) *InviteLink {
	return &InviteLink{
		referralService: referralService,
		notifier:        notifier,
	}
}

func (l *InviteLink) SendInvite(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	code, err := l.referralService.CodeByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("code by player_id: %w", err)
	}

	if err := l.notifier.SendJoinReferralLink(
		ctx,
		chatID,
		code,
	); err != nil {
		return fmt.Errorf("send invite link: %w", err)
	}

	return nil
}
