package postgres_test

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
)

func (s *PostgresSuit) TestInsertReferral() {
	s.Run("insert referral with invited_by_chat_id", func() {
		var (
			ctx       = context.Background()
			chatID    = generateChatID()
			invitedBy = generateChatID()
			createdAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			code      = referral.CodeFromString("referral-code-insert-1")
			ref       = referral.Referral{
				ChatID:    chatID,
				InvitedBy: invitedBy,
				Code:      code,
				CreatedAt: createdAt,
			}
		)

		err := s.pgDB.InsertReferral(ctx, ref)
		s.Require().NoError(err)

		dbRef := s.getReferralByChatID(ctx, chatID)

		s.compareReferrals(ref, dbRef)
		s.Require().Equal(1, s.countReferrals(ctx))
	})

	s.Run("insert referral without invited_by_chat_id stores null", func() {
		var (
			ctx       = context.Background()
			chatID    = generateChatID()
			createdAt = time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
			code      = referral.CodeFromString("referral-code-insert-2")
			ref       = referral.Referral{
				ChatID:    chatID,
				InvitedBy: msginfo.ChatIDFromInt64(0),
				Code:      code,
				CreatedAt: createdAt,
			}
		)

		err := s.pgDB.InsertReferral(ctx, ref)
		s.Require().NoError(err)

		dbRef := s.getReferralByChatID(ctx, chatID)

		s.compareReferrals(ref, dbRef)
		s.Require().Equal(1, s.countReferrals(ctx))
	})

	s.Run("insert referral with duplicate chat_id returns constraint error", func() {
		var (
			ctx       = context.Background()
			chatID    = generateChatID()
			invitedBy = generateChatID()
			createdAt = time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
			code1     = referral.CodeFromString("referral-code-insert-3a")
			code2     = referral.CodeFromString("referral-code-insert-3b")
		)

		err := s.pgDB.InsertReferral(ctx, referral.Referral{
			ChatID:    chatID,
			InvitedBy: invitedBy,
			Code:      code1,
			CreatedAt: createdAt,
		})
		s.Require().NoError(err)

		// A second insert for the same chat_id must violate the primary key.
		err = s.pgDB.InsertReferral(ctx, referral.Referral{
			ChatID:    chatID,
			InvitedBy: invitedBy,
			Code:      code2,
			CreatedAt: createdAt.Add(time.Hour),
		})
		s.Require().Error(err)
		s.Require().True(driver.NewPgx().IsConstraintError(err, "referral_pkey"))

		s.Require().Equal(1, s.countReferrals(ctx))
	})

	s.Run("insert referral with duplicate code returns constraint error", func() {
		var (
			ctx       = context.Background()
			chatID1   = generateChatID()
			chatID2   = generateChatID()
			invitedBy = generateChatID()
			createdAt = time.Date(2026, 4, 5, 6, 7, 8, 0, time.UTC)
			code      = referral.CodeFromString("referral-code-insert-4")
		)

		err := s.pgDB.InsertReferral(ctx, referral.Referral{
			ChatID:    chatID1,
			InvitedBy: invitedBy,
			Code:      code,
			CreatedAt: createdAt,
		})
		s.Require().NoError(err)

		// A second insert reusing the same code must violate the unique index on code.
		err = s.pgDB.InsertReferral(ctx, referral.Referral{
			ChatID:    chatID2,
			InvitedBy: invitedBy,
			Code:      code,
			CreatedAt: createdAt.Add(time.Hour),
		})
		s.Require().Error(err)
		s.Require().True(driver.NewPgx().IsConstraintError(err, "referral_code_idx"))

		s.Require().Equal(1, s.countReferrals(ctx))
	})

	s.Run("insert multiple referrals with distinct codes and chat ids", func() {
		var (
			ctx       = context.Background()
			chatID1   = generateChatID()
			chatID2   = generateChatID()
			chatID3   = generateChatID()
			invitedBy = generateChatID()
			createdAt = time.Date(2026, 5, 6, 7, 8, 9, 0, time.UTC)
			code1     = referral.CodeFromString("referral-code-insert-5a")
			code2     = referral.CodeFromString("referral-code-insert-5b")
			code3     = referral.CodeFromString("referral-code-insert-5c")

			ref1 = referral.Referral{
				ChatID:    chatID1,
				InvitedBy: invitedBy,
				Code:      code1,
				CreatedAt: createdAt,
			}

			ref2 = referral.Referral{
				ChatID:    chatID2,
				InvitedBy: msginfo.ChatIDFromInt64(0),
				Code:      code2,
				CreatedAt: createdAt.Add(time.Hour),
			}

			ref3 = referral.Referral{
				ChatID:    chatID3,
				InvitedBy: invitedBy,
				Code:      code3,
				CreatedAt: createdAt.Add(2 * time.Hour),
			}
		)

		err := s.pgDB.InsertReferral(ctx, ref1)
		s.Require().NoError(err)

		err = s.pgDB.InsertReferral(ctx, ref2)
		s.Require().NoError(err)

		err = s.pgDB.InsertReferral(ctx, ref3)
		s.Require().NoError(err)

		s.Require().Equal(3, s.countReferrals(ctx))

		dbRef1 := s.getReferralByChatID(ctx, chatID1)
		s.compareReferrals(ref1, dbRef1)

		dbRef2 := s.getReferralByChatID(ctx, chatID2)
		s.compareReferrals(ref2, dbRef2)

		dbRef3 := s.getReferralByChatID(ctx, chatID3)
		s.compareReferrals(ref3, dbRef3)
	})
}

func (s *PostgresSuit) compareReferrals(expected, actual referral.Referral) {
	expected.CreatedAt = expected.CreatedAt.UTC()
	actual.CreatedAt = actual.CreatedAt.UTC()
	s.Require().Equal(expected, actual)
}

func (s *PostgresSuit) getReferralByChatID(ctx context.Context, chatID msginfo.ChatID) referral.Referral {
	ref, err := s.pgDB.GetReferralByChatID(ctx, chatID)
	s.Require().NoError(err)

	return ref
}

func (s *PostgresSuit) countReferrals(ctx context.Context) int {
	var (
		query = `
			SELECT COUNT(*)
			FROM referral
		`
		trx   = s.pgDB.Transactor()
		count int
	)

	err := sqlx.GetContext(ctx, trx.ExtContext(ctx), &count, query)
	s.Require().NoError(err)

	return count
}
