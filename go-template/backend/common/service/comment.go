package service

import (
	"github.com/uptrace/bun"

	"github.com/hayrat/go-template2/backend/common/model"
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
