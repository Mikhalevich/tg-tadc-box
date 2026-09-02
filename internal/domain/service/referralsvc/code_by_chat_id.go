package referralsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

func (s *Service) CodeByChatID(
	ctx context.Context,
	chatID msginfo.ChatID,
) (referral.Code, error) {
	ref, err := s.getOrCreateReferral(ctx, chatID)
	if err != nil {
		return referral.EmptyCode(), fmt.Errorf("get or create referral: %w", err)
	}

	return ref.Code, nil
}

func (s *Service) getOrCreateReferral(
	ctx context.Context,
	chatID msginfo.ChatID,
) (referral.Referral, error) {
	ref, err := s.repo.GetReferralByChatID(ctx, chatID)
	if err != nil {
		if !perror.IsType(err, perror.TypeNotFound) {
			return referral.Referral{}, fmt.Errorf("get referral from repo: %w", err)
		}

		newRef := referral.Referral{
			ChatID: chatID,
			Code:   referral.GenerateCode(),
		}

		if err := s.repo.InsertReferral(ctx, newRef); err != nil {
			return referral.Referral{}, fmt.Errorf("insert referral: %w", err)
		}

		return newRef, nil
	}

	return ref, nil
}
