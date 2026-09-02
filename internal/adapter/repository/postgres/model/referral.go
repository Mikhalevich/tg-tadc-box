package model

import (
	"database/sql"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

type Referral struct {
	ChatID    int64         `db:"chat_id"`
	InvitedBy sql.NullInt64 `db:"invited_by_chat_id"`
	Code      string        `db:"code"`
	CreatedAt time.Time     `db:"created_at"`
}

func (r Referral) ToDom() referral.Referral {
	return referral.Referral{
		ChatID:    msginfo.ChatIDFromInt64(r.ChatID),
		InvitedBy: msginfo.ChatIDFromInt64(r.InvitedBy.Int64),
		Code:      referral.CodeFromString(r.Code),
		CreatedAt: r.CreatedAt,
	}
}

func ToDBReferral(dom referral.Referral) Referral {
	return Referral{
		ChatID:    dom.ChatID.Int64(),
		InvitedBy: convertChatIDToNullInt64(dom.InvitedBy),
		Code:      dom.Code.String(),
		CreatedAt: dom.CreatedAt,
	}
}

func convertChatIDToNullInt64(chatID msginfo.ChatID) sql.NullInt64 {
	return sql.NullInt64{
		Int64: chatID.Int64(),
		Valid: chatID.Int64() != 0,
	}
}
