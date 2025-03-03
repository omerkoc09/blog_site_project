package service

import (
	"github.com/uptrace/bun"

	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	"github.com/hayrat/go-template2/backend/pkg/service"
)

type IPostService interface {
	service.IBaseService[model.Post]
	DeleteImage(ctx *app.Ctx) error
}

type PostService struct {
	service.BaseService[model.Post]
}

func NewPostService(db *bun.DB) IPostService {
	return &PostService{
		service.BaseService[model.Post]{DB: db},
	}
}

func (s PostService) DeleteImage(ctx *app.Ctx) error {
	postID := ctx.ParamsInt64("id")
	if postID == 0 {
		return errorsx.BadRequestError("Invalid post ID")
	}

	// Post'u veritabanından al
	m, err := s.GetByID(ctx.Context(), postID)
	if err != nil {
		return err
	}

	// Image alanını boşalt
	m.Image = ""

	// Post'u güncelle
	err = s.Update(ctx.Context(), m)
	if err != nil {
		return err
	}

	return nil
}
