package model

import (
	"database/sql"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type Box struct {
	ID                  int          `db:"id"`
	ChatID              int64        `db:"chat_id"`
	Status              string       `db:"status"`
	Type                string       `db:"type"`
	CreatedAt           time.Time    `db:"created_at"`
	AvailableAt         time.Time    `db:"available_at"`
	ReadyNotificationAt sql.NullTime `db:"ready_notification_at"`
	CompletedAt         sql.NullTime `db:"completed_at"`
}

func (b Box) ToDomBox() box.Box {
	return box.Box{
		ID:                  b.ID,
		ChatID:              msginfo.ChatIDFromInt64(b.ChatID),
		Status:              box.Status(b.Status),
		Type:                box.Type(b.Type),
		CreatedAt:           b.CreatedAt,
		AvailableAt:         b.AvailableAt,
		ReadyNotificationAt: b.ReadyNotificationAt.Time,
		CompletedAt:         b.CompletedAt.Time,
	}
}

func ToDomBoxes(dbBoxes []Box) []box.Box {
	if len(dbBoxes) == 0 {
		return nil
	}

	domBoxes := make([]box.Box, 0, len(dbBoxes))
	for _, dbBox := range dbBoxes {
		domBoxes = append(domBoxes, dbBox.ToDomBox())
	}

	return domBoxes
}

func ToDBBox(domBox box.Box) Box {
	return Box{
		ID:                  domBox.ID,
		ChatID:              domBox.ChatID.Int64(),
		Status:              domBox.Status.String(),
		Type:                domBox.Type.String(),
		CreatedAt:           domBox.CreatedAt,
		AvailableAt:         domBox.AvailableAt,
		ReadyNotificationAt: toNullTime(domBox.ReadyNotificationAt),
		CompletedAt:         toNullTime(domBox.CompletedAt),
	}
}

func toNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}
