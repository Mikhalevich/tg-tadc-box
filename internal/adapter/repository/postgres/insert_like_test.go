package postgres_test

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/driver"
	"github.com/Mikhalevich/tg-tadc-box/internal/adapter/repository/postgres/model"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/like"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (s *PostgresSuit) TestInsertLike() {
	s.Run("insert like", func() {
		var (
			ctx       = context.Background()
			chatID    = generateChatID()
			createdAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			boxID     = s.insertBox(ctx, chatID)
			rwdID     = s.insertReward(ctx, reward.RewardTypeCommon, "https://example.com/common.png", createdAt)
		)
		err := s.pgDB.InsertLike(ctx, like.Like{
			BoxID:     boxID,
			ChatID:    chatID,
			RewardID:  rwdID,
			Type:      like.TypeLike,
			CreatedAt: createdAt,
		})
		s.Require().NoError(err)

		dbLike := s.getLikeByBoxID(ctx, boxID)
		s.Require().Equal(boxID.Int(), dbLike.BoxID)
		s.Require().Equal(chatID, dbLike.ChatID)
		s.Require().Equal(rwdID.Int(), dbLike.RewardID)
		s.Require().Equal(like.TypeLike.String(), dbLike.Type)
		s.Require().WithinDuration(createdAt, dbLike.CreatedAt, time.Second)
		s.Require().Equal(1, s.countLikes(ctx))
	})

	s.Run("insert dislike", func() {
		var (
			ctx       = context.Background()
			chatID    = generateChatID()
			createdAt = time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
			boxID     = s.insertBox(ctx, chatID)
			rwdID     = s.insertReward(ctx, reward.RewardTypeRare, "https://example.com/rare.png", createdAt)
		)

		err := s.pgDB.InsertLike(ctx, like.Like{
			BoxID:     boxID,
			ChatID:    chatID,
			RewardID:  rwdID,
			Type:      like.TypeDislike,
			CreatedAt: createdAt,
		})
		s.Require().NoError(err)

		dbLike := s.getLikeByBoxID(ctx, boxID)
		s.Require().Equal(like.TypeDislike.String(), dbLike.Type)
		s.Require().Equal(1, s.countLikes(ctx))
	})

	s.Run("duplicate like for the same box returns constraint error", func() {
		var (
			ctx       = context.Background()
			chatID    = generateChatID()
			createdAt = time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
			boxID     = s.insertBox(ctx, chatID)
			rwdID     = s.insertReward(ctx, reward.RewardTypeEpic, "https://example.com/epic.png", createdAt)
		)

		err := s.pgDB.InsertLike(ctx, like.Like{
			BoxID:     boxID,
			ChatID:    chatID,
			RewardID:  rwdID,
			Type:      like.TypeLike,
			CreatedAt: createdAt,
		})
		s.Require().NoError(err)

		// A second like/dislike for the same box must violate the unique index.
		err = s.pgDB.InsertLike(ctx, like.Like{
			BoxID:     boxID,
			ChatID:    chatID,
			RewardID:  rwdID,
			Type:      like.TypeDislike,
			CreatedAt: createdAt.Add(time.Hour),
		})
		s.Require().Error(err)
		s.Require().True(driver.NewPgx().IsConstraintError(err, "likes_box_id_idx"))

		s.Require().Equal(1, s.countLikes(ctx))
	})

	s.Run("insert like with non-existent box returns fk error", func() {
		var (
			ctx    = context.Background()
			chatID = generateChatID()
			rwdID  = s.insertReward(ctx, reward.RewardTypeLegendary, "https://example.com/legendary.png", time.Now())
		)

		err := s.pgDB.InsertLike(ctx, like.Like{
			BoxID:     box.IDFromInt(999999),
			ChatID:    chatID,
			RewardID:  rwdID,
			Type:      like.TypeLike,
			CreatedAt: time.Now(),
		})
		s.Require().Error(err)
		s.Require().True(driver.NewPgx().IsConstraintError(err, "likes_box_id_fk"))

		s.Require().Equal(0, s.countLikes(ctx))
	})

	s.Run("insert like with non-existent reward returns fk error", func() {
		var (
			ctx    = context.Background()
			chatID = generateChatID()
			boxID  = s.insertBox(ctx, chatID)
		)

		err := s.pgDB.InsertLike(ctx, like.Like{
			BoxID:     boxID,
			ChatID:    chatID,
			RewardID:  reward.IDFromInt(999999),
			Type:      like.TypeLike,
			CreatedAt: time.Now(),
		})
		s.Require().Error(err)
		s.Require().True(driver.NewPgx().IsConstraintError(err, "likes_reward_id_fk"))

		s.Require().Equal(0, s.countLikes(ctx))
	})
}

func (s *PostgresSuit) insertBox(ctx context.Context, chatID msginfo.ChatID) box.ID {
	now := time.Now()

	boxID, err := s.pgDB.InsertBox(ctx, box.Box{
		ChatID:      chatID,
		Status:      box.StatusPending,
		Type:        box.TypeCommon,
		CreatedAt:   now,
		AvailableAt: now,
	})
	s.Require().NoError(err)

	return box.IDFromInt(boxID)
}

func (s *PostgresSuit) insertReward(
	ctx context.Context,
	rwdType reward.RewardType,
	uri string,
	createdAt time.Time,
) reward.ID {
	var (
		query = `
			INSERT INTO reward(
				type,
				uri,
				created_at
			) VALUES (
				$1,
				$2,
				$3
			)
			RETURNING id
		`
		trx      = s.pgDB.Transactor()
		rewardID int
	)

	err := sqlx.GetContext(ctx, trx.ExtContext(ctx), &rewardID, query,
		rwdType.String(),
		uri,
		createdAt,
	)
	s.Require().NoError(err)

	return reward.IDFromInt(rewardID)
}

func (s *PostgresSuit) getLikeByBoxID(ctx context.Context, boxID box.ID) model.Like {
	var (
		query = `
			SELECT
				id,
				box_id,
				chat_id,
				reward_id,
				type,
				created_at
			FROM
				likes
			WHERE
				box_id = $1
		`
		trx    = s.pgDB.Transactor()
		dbLike model.Like
	)

	err := sqlx.GetContext(ctx, trx.ExtContext(ctx), &dbLike, query, boxID.Int())
	s.Require().NoError(err)

	return dbLike
}

func (s *PostgresSuit) countLikes(ctx context.Context) int {
	var (
		query = `
			SELECT COUNT(*)
			FROM likes
		`
		trx   = s.pgDB.Transactor()
		count int
	)

	err := sqlx.GetContext(ctx, trx.ExtContext(ctx), &count, query)
	s.Require().NoError(err)

	return count
}

func generateChatID() msginfo.ChatID {
	//nolint:gosec
	return msginfo.ChatIDFromInt64(rand.Int64N(1000000))
}
