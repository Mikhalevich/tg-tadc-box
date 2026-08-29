package postgres_test

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (s *PostgresSuit) TestCompleteBox() {
	s.Run("complete box from pending to opened", func() {
		var (
			ctx         = context.Background()
			chatID      = generateChatID()
			createdAt   = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			completedAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			boxID       = s.insertBoxWithStatus(ctx, chatID, createdAt, box.StatusPending)
		)

		err := s.pgDB.CompleteBox(ctx, boxID, box.StatusOpened, completedAt)
		s.Require().NoError(err)

		dbBox := s.getBoxByID(ctx, boxID)
		s.Require().Equal(box.StatusOpened.String(), dbBox.Status)
		s.Require().True(dbBox.CompletedAt.Valid)
	})

	s.Run("complete box from in_progress to opened", func() {
		var (
			ctx         = context.Background()
			chatID      = generateChatID()
			createdAt   = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			completedAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			boxID       = s.insertBoxWithStatus(ctx, chatID, createdAt, box.StatusInProgress)
		)

		err := s.pgDB.CompleteBox(ctx, boxID, box.StatusOpened, completedAt)
		s.Require().NoError(err)

		dbBox := s.getBoxByID(ctx, boxID)
		s.Require().Equal(box.StatusOpened.String(), dbBox.Status)
		s.Require().True(dbBox.CompletedAt.Valid)
	})

	s.Run("complete box from pending to canceled", func() {
		var (
			ctx         = context.Background()
			chatID      = generateChatID()
			createdAt   = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			completedAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			boxID       = s.insertBoxWithStatus(ctx, chatID, createdAt, box.StatusPending)
		)

		err := s.pgDB.CompleteBox(ctx, boxID, box.StatusCanceled, completedAt)
		s.Require().NoError(err)

		dbBox := s.getBoxByID(ctx, boxID)
		s.Require().Equal(box.StatusCanceled.String(), dbBox.Status)
		s.Require().True(dbBox.CompletedAt.Valid)
	})

	s.Run("complete box from opened to canceled", func() {
		var (
			ctx         = context.Background()
			chatID      = generateChatID()
			createdAt   = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			completedAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			boxID       = s.insertBoxWithStatus(ctx, chatID, createdAt, box.StatusOpened)
		)

		err := s.pgDB.CompleteBox(ctx, boxID, box.StatusCanceled, completedAt)
		s.Require().NoError(err)

		dbBox := s.getBoxByID(ctx, boxID)
		s.Require().Equal(box.StatusCanceled.String(), dbBox.Status)
		s.Require().True(dbBox.CompletedAt.Valid)
	})

	s.Run("complete non-existent box returns no rows updated", func() {
		var (
			ctx         = context.Background()
			completedAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			boxID       = box.IDFromInt(999999)
		)

		err := s.pgDB.CompleteBox(ctx, boxID, box.StatusOpened, completedAt)
		s.Require().ErrorIs(err, perror.NoRowsUpdated())
	})
}

func (s *PostgresSuit) insertBoxWithStatus(
	ctx context.Context,
	chatID msginfo.ChatID,
	createdAt time.Time,
	status box.Status,
) box.ID {
	boxID, err := s.pgDB.InsertBox(ctx, box.Box{
		ChatID:      chatID,
		Status:      status,
		Type:        box.TypeCommon,
		CreatedAt:   createdAt,
		AvailableAt: createdAt,
	})
	s.Require().NoError(err)

	return box.IDFromInt(boxID)
}

func (s *PostgresSuit) getBoxByID(ctx context.Context, boxID box.ID) model.Box {
	var (
		query = `
			SELECT
				id,
				chat_id,
				status,
				type,
				created_at,
				available_at,
				ready_notification_at,
				completed_at,
				meta
			FROM
				box
			WHERE
				id = $1
		`
		trx   = s.pgDB.Transactor()
		dbBox model.Box
	)

	err := sqlx.GetContext(ctx, trx.ExtContext(ctx), &dbBox, query, boxID.Int())
	s.Require().NoError(err)

	return dbBox
}
