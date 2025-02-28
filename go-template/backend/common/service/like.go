package service

import (
	"context"

	"github.com/uptrace/bun"

	"database/sql"

	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	"github.com/hayrat/go-template2/backend/pkg/service"
)

type ILikeService interface {
	service.IBaseService[model.Like]
	GetLikedByUser(ctx context.Context, userID, postID int64) (bool, *model.Like, error)
}

type LikeService struct {
	service.BaseService[model.Like]
}

func NewLikeService(db *bun.DB) ILikeService {
	return &LikeService{
		service.BaseService[model.Like]{DB: db},
	}
}

func (s LikeService) Create(ctx context.Context, m *model.Like) error {

	_, err := s.DB.NewInsert().Model(m).Exec(ctx)
	return errorsx.Database(err)
}

func (s LikeService) Update(ctx context.Context, m model.Like) error {

	_, err := s.DB.NewUpdate().
		Model(&m).
		WherePK().
		Exec(ctx)
	return errorsx.Database(err)
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
