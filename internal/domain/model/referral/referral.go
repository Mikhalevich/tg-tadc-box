package referral

import (
	"time"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type Code string

func (c Code) String() string {
	return string(c)
}

func CodeFromString(code string) Code {
	return Code(code)
}

func EmptyCode() Code {
	return Code("")
}

func GenerateCode() Code {
	return Code(uuid.NewString())
}

type Referral struct {
	ChatID    msginfo.ChatID
	InvitedBy msginfo.ChatID
	Code      Code
	CreatedAt time.Time
}
