package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
)

func (m *MessageProcessor) SetButtonRows(
	ctx context.Context,
	rows ...button.ButtonRow,
) ([]button.InlineKeyboardButtonRow, error) {
	if len(rows) == 0 {
		return nil, nil
	}

	inlineButtonRows := make([]button.InlineKeyboardButtonRow, 0, len(rows))

	for _, row := range rows {
		buttonRow := make([]button.InlineKeyboardButton, 0, len(row))

		for _, btn := range row {
			buttonRow = append(buttonRow, button.InlineKeyboardButton{
				ID:      btn.ID,
				Caption: btn.Caption,
			})
		}

		if len(buttonRow) > 0 {
			inlineButtonRows = append(inlineButtonRows, buttonRow)
		}
	}

	rows = filterNoOperation(rows)

	if len(rows) > 0 {
		if err := m.buttonRepository.SetButtonRows(ctx, rows...); err != nil {
			return nil, fmt.Errorf("set button rows: %w", err)
		}
	}

	return inlineButtonRows, nil
}

func filterNoOperation(rows []button.ButtonRow) []button.ButtonRow {
	for rowIdx := len(rows) - 1; rowIdx >= 0; rowIdx-- {
		var (
			row    = rows[rowIdx]
			rowLen = len(row)
		)

		for elemIdx := len(row) - 1; elemIdx >= 0; elemIdx-- {
			if row[elemIdx].Operation.IsNoOperation() {
				row = append(row[0:elemIdx], row[elemIdx+1:]...)
			}
		}

		if len(row) == rowLen {
			continue
		}

		if len(row) == 0 {
			rows = append(rows[:rowIdx], rows[rowIdx+1:]...)
		} else {
			rows[rowIdx] = row
		}
	}

	return rows
}
