package messagesender

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Sender interface {
	SendMessage(
		ctx context.Context,
		msg msginfo.Message,
	) error
}

type ImageProvider interface {
	Reward(
		ctx context.Context,
		rew reward.Reward,
	) ([]byte, error)
	Chest(ctx context.Context, boxType box.Type) ([]byte, error)
}

type MessageSender struct {
	sender        Sender
	imageProvider ImageProvider
}

func New(
	sender Sender,
	imageProvider ImageProvider,
) *MessageSender {
	return &MessageSender{
		sender:        sender,
		imageProvider: imageProvider,
	}
}
