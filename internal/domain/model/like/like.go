package like

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Type string

const (
	TypeLike    Type = "like"
	TypeDislike Type = "dislike"
)

func (t Type) String() string {
	return string(t)
}

type Like struct {
	ID        int
	BoxID     box.ID
	ChatID    msginfo.ChatID
	RewardID  reward.ID
	Type      Type
	CreatedAt time.Time
}
