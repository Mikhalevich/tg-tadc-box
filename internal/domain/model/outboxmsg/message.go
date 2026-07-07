package outboxmsg

import (
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusDispatched Status = "dispatched"
	StatusCanceled   Status = "canceled"
)

type Message struct {
	msginfo.Message

	ID         int
	RetryCount int
}
