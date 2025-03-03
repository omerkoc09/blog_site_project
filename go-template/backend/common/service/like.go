package service

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"

	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	"github.com/hayrat/go-template2/backend/pkg/service"
)

type ILikeService interface {
	service.IBaseService[model.Like]
	GetLikedByUser(ctx context.Context, userID, postID int64) (bool, *model.Like, error)
	ToggleLike(ctx context.Context, userID, postID int64) error
}

type LikeService struct {
	service.BaseService[model.Like]
}

func NewLikeService(db *bun.DB) ILikeService {
	return &LikeService{
		service.BaseService[model.Like]{DB: db},
	}
}

func (s LikeService) GetLikedByUser(ctx context.Context, userID, postID int64) (bool, *model.Like, error) {
	like := new(model.Like)
	err := s.DB.NewSelect().
		Model(like).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, nil
		}
		return false, nil, errorsx.Database(err)
	}
	return true, like, nil
}

func (s LikeService) ToggleLike(ctx context.Context, userID, postID int64) error {
	isLiked, like, err := s.GetLikedByUser(ctx, userID, postID)
	if err != nil {
		return err
	}

	if isLiked {
		return s.Delete(ctx, like.ID)
	}

	like = &model.Like{
		PostId: postID,
		UserId: userID,
	}
	return s.Create(ctx, like)
}
