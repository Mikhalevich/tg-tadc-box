package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func (m *MessageProcessor) DeleteMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	if messageID.Int() == 0 {
		return nil
	}

	if err := m.sender.DeleteMessage(ctx, chatID, messageID); err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}
