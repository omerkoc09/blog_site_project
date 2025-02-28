package service

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	"github.com/hayrat/go-template2/backend/pkg/service"
)

type ICommentService interface {
	service.IBaseService[model.Comment]
	//GetByEmail(ctx context.Context, email string) (model.Comment, error)
}

type CommentService struct {
	service.BaseService[model.Comment]
}

func NewCommentService(db *bun.DB) ICommentService {
	return &CommentService{
		service.BaseService[model.Comment]{DB: db},
	}
}

func (s CommentService) Create(ctx context.Context, m *model.Comment) error {

	_, err := s.DB.NewInsert().Model(m).Exec(ctx)
	return errorsx.Database(err)
}

func (s CommentService) Update(ctx context.Context, m model.Comment) error {

	_, err := s.DB.NewUpdate().
		Model(&m).
		WherePK().
		Exec(ctx)
	return errorsx.Database(err)
}
