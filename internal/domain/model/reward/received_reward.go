package reward

import (
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type ReceivedReward struct {
	ID        ID
	ChatID    msginfo.ChatID
	RewardID  ID
	BoxID     box.ID
	CreatedAt time.Time
}
