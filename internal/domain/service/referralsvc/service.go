package referralsvc

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

type Repository interface {
	GetReferralByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (referral.Referral, error)
	InsertReferral(
		ctx context.Context,
		ref referral.Referral,
	) error
}

type Service struct {
	repo Repository
}

func New(
	repo Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}
