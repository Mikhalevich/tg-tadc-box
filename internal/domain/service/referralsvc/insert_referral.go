package referralsvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

func (s *Service) InsertReferral(
	ctx context.Context,
	ref referral.Referral,
) error {
	if err := s.repo.InsertReferral(ctx, ref); err != nil {
		return fmt.Errorf("insert referral: %w", err)
	}

	return nil
}
