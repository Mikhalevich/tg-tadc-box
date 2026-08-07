package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/internal/jsonb"
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
	Payload             jsonb.JSONB  `db:"meta"`
}

func (b Box) ToDomBox() (box.Box, error) {
	var meta box.Meta
	if err := jsonb.ConvertTo(b.Payload, &meta); err != nil {
		return box.Box{}, fmt.Errorf("convert to box meta: %w", err)
	}

	return box.Box{
		ID:                  box.IDFromInt(b.ID),
		ChatID:              msginfo.ChatIDFromInt64(b.ChatID),
		Status:              box.Status(b.Status),
		Type:                box.Type(b.Type),
		CreatedAt:           b.CreatedAt,
		AvailableAt:         b.AvailableAt,
		ReadyNotificationAt: b.ReadyNotificationAt.Time,
		CompletedAt:         b.CompletedAt.Time,
		Meta:                meta,
	}, nil
}

func ToDomBoxes(dbBoxes []Box) ([]box.Box, error) {
	if len(dbBoxes) == 0 {
		return nil, nil
	}

	domBoxes := make([]box.Box, 0, len(dbBoxes))
	for _, dbBox := range dbBoxes {
		domBox, err := dbBox.ToDomBox()
		if err != nil {
			return nil, fmt.Errorf("convert to dom box: %w", err)
		}
		domBoxes = append(domBoxes, domBox)
	}

	return domBoxes, nil
}

func ToDBBox(domBox box.Box) (Box, error) {
	payload, err := jsonb.NewFromMarshaler(domBox.Meta)
	if err != nil {
		return Box{}, fmt.Errorf("marshal meta to jsonb: %w", err)
	}

	return Box{
		ID:                  domBox.ID.Int(),
		ChatID:              domBox.ChatID.Int64(),
		Status:              domBox.Status.String(),
		Type:                domBox.Type.String(),
		CreatedAt:           domBox.CreatedAt,
		AvailableAt:         domBox.AvailableAt,
		ReadyNotificationAt: toNullTime(domBox.ReadyNotificationAt),
		CompletedAt:         toNullTime(domBox.CompletedAt),
		Payload:             payload,
	}, nil
}

func toNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}
