package referralsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

func (s *Service) GetReferralByCode(
	ctx context.Context,
	code referral.Code,
) (referral.Referral, error) {
	ref, err := s.repo.GetReferralByCode(ctx, code)
	if err != nil {
		return referral.Referral{}, fmt.Errorf("get referral by code: %w", err)
	}

	return ref, nil
}
