package postgres_test

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

func (s *PostgresSuit) TestGetReferralByChatID() {
	s.Run("get existing referral with invited_by_chat_id", func() {
		var (
			ctx          = context.Background()
			chatID       = generateChatID()
			invitedBy    = generateChatID()
			createdAt    = time.Date(2026, 4, 5, 6, 7, 8, 0, time.UTC)
			code         = referral.CodeFromString("referral-code-1")
			expectedChat = msginfo.ChatIDFromInt64(chatID.Int64())
		)

		s.insertReferral(ctx, expectedChat, invitedBy, code, createdAt)

		ref, err := s.pgDB.GetReferralByChatID(ctx, expectedChat)
		s.Require().NoError(err)
		s.Require().Equal(expectedChat.Int64(), ref.ChatID.Int64())
		s.Require().Equal(invitedBy.Int64(), ref.InvitedBy.Int64())
		s.Require().Equal(code.String(), ref.Code.String())
		s.Require().WithinDuration(createdAt, ref.CreatedAt, time.Second)
	})

	s.Run("get existing referral without invited_by_chat_id", func() {
		var (
			ctx       = context.Background()
			chatID    = generateChatID()
			createdAt = time.Date(2026, 5, 6, 7, 8, 9, 0, time.UTC)
			code      = referral.CodeFromString("referral-code-2")
		)

		s.insertReferral(ctx, chatID, msginfo.ChatIDFromInt64(0), code, createdAt)

		ref, err := s.pgDB.GetReferralByChatID(ctx, chatID)
		s.Require().NoError(err)
		s.Require().Equal(chatID.Int64(), ref.ChatID.Int64())
		s.Require().Zero(ref.InvitedBy.Int64())
		s.Require().Equal(code.String(), ref.Code.String())
		s.Require().WithinDuration(createdAt, ref.CreatedAt, time.Second)
	})

	s.Run("get non-existent referral returns not found", func() {
		var (
			ctx = context.Background()
			//nolint:gosec
			chatID = msginfo.ChatIDFromInt64(rand.Int64N(1000000) + 1000000)
		)

		ref, err := s.pgDB.GetReferralByChatID(ctx, chatID)
		s.Require().Error(err)
		s.Require().True(perror.IsType(err, perror.TypeNotFound))
		s.Require().Zero(ref)
	})

	s.Run("get referral returns the correct one when multiple exist", func() {
		var (
			ctx       = context.Background()
			chatID1   = generateChatID()
			chatID2   = generateChatID()
			invitedBy = generateChatID()
			createdAt = time.Date(2026, 6, 7, 8, 9, 10, 0, time.UTC)
			code1     = referral.CodeFromString("referral-code-3")
			code2     = referral.CodeFromString("referral-code-4")
		)

		s.insertReferral(ctx, chatID1, invitedBy, code1, createdAt)
		s.insertReferral(ctx, chatID2, invitedBy, code2, createdAt.Add(time.Hour))

		ref, err := s.pgDB.GetReferralByChatID(ctx, chatID1)
		s.Require().NoError(err)
		s.Require().Equal(chatID1.Int64(), ref.ChatID.Int64())
		s.Require().Equal(invitedBy.Int64(), ref.InvitedBy.Int64())
		s.Require().Equal(code1.String(), ref.Code.String())
		s.Require().WithinDuration(createdAt, ref.CreatedAt, time.Second)

		ref, err = s.pgDB.GetReferralByChatID(ctx, chatID2)
		s.Require().NoError(err)
		s.Require().Equal(chatID2.Int64(), ref.ChatID.Int64())
		s.Require().Equal(invitedBy.Int64(), ref.InvitedBy.Int64())
		s.Require().Equal(code2.String(), ref.Code.String())
		s.Require().WithinDuration(createdAt.Add(time.Hour), ref.CreatedAt, time.Second)
	})
}

func (s *PostgresSuit) insertReferral(
	ctx context.Context,
	chatID msginfo.ChatID,
	invitedBy msginfo.ChatID,
	code referral.Code,
	createdAt time.Time,
) {
	var (
		query = `
			INSERT INTO referral(
				chat_id,
				invited_by_chat_id,
				code,
				created_at
			) VALUES (
				$1,
				$2,
				$3,
				$4
			)
		`
		trx = s.pgDB.Transactor()
	)

	invitedByArg := any(nil)
	if invitedBy.Int64() != 0 {
		invitedByArg = invitedBy.Int64()
	}

	_, err := trx.ExtContext(ctx).ExecContext(ctx, query,
		chatID.Int64(),
		invitedByArg,
		code.String(),
		createdAt,
	)
	s.Require().NoError(err)
}
