package model

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
)

type Player struct {
	ID               int         `db:"id"`
	ChatID           int64       `db:"chat_id"`
	CreatedAt        time.Time   `db:"created_at"`
	Payload          jsonb.JSONB `db:"payload"`
	PayloadVersion   int         `db:"payload_version"`
	PayloadUpdatedAt time.Time   `db:"payload_updated_at"`
}

func (p Player) ToDom() (player.Player, error) {
	var profile player.Profile
	if err := jsonb.ConvertTo(p.Payload, &profile); err != nil {
		return player.Player{}, fmt.Errorf("convert to user profile: %w", err)
	}

	return player.Player{
		ID:               player.IDFromInt(p.ID),
		ChatID:           msginfo.ChatIDFromInt64(p.ChatID),
		CreatedAt:        p.CreatedAt,
		Profile:          profile,
		ProfileVersion:   p.PayloadVersion,
		ProfileUpdatedAt: p.PayloadUpdatedAt,
	}, nil
}

func ToDBPlayer(domPlayer player.Player) (Player, error) {
	payload, err := jsonb.NewFromMarshaler(domPlayer.Profile)
	if err != nil {
		return Player{}, fmt.Errorf("convert profile to jsonb: %w", err)
	}

	return Player{
		ID:               domPlayer.ID.Int(),
		ChatID:           domPlayer.ChatID.Int64(),
		CreatedAt:        domPlayer.CreatedAt,
		Payload:          payload,
		PayloadVersion:   domPlayer.ProfileVersion,
		PayloadUpdatedAt: domPlayer.ProfileUpdatedAt,
	}, nil
}
