package messagesender

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/outboxmsg/imagepayload"
)

func (m *MessageSender) SendMessage(
	ctx context.Context,
	msg msginfo.Message,
) error {
	if msg.Type == msginfo.MessageTypePNG ||
		msg.Type == msginfo.MessageTypeEditPNG {
		payload, err := m.loadPayload(ctx, msg.Payload)
		if err != nil {
			return fmt.Errorf("load payload: %w", err)
		}

		msg.Payload = payload
	}

	if err := m.sender.SendMessage(ctx, msg); err != nil {
		return fmt.Errorf("sender send message: %w", err)
	}

	return nil
}

func (m *MessageSender) loadPayload(ctx context.Context, payload []byte) ([]byte, error) {
	imgPayload, err := imagepayload.GOBDecode(payload)
	if err != nil {
		return nil, fmt.Errorf("decode image payload: %w", err)
	}

	switch imgPayload.Type {
	case imagepayload.PayloadTypeCommonChest:
		image, err := m.imageProvider.Chest(ctx)
		if err != nil {
			return nil, fmt.Errorf("get chest: %w", err)
		}

		return image, nil

	case imagepayload.PayloadTypeReward:
		image, err := m.imageProvider.Reward(ctx, imgPayload.Reward)
		if err != nil {
			return nil, fmt.Errorf("get reward: %w", err)
		}

		return image, nil
	}

	return nil, fmt.Errorf("invaliad type %v", imgPayload.Type)
}
