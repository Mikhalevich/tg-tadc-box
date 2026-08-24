package likesvc

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/like"
)

type Repository interface {
	InsertLike(
		ctx context.Context,
		rwdLike like.Like,
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
