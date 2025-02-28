package handlers

import (
	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/common/service"
	"github.com/hayrat/go-template2/backend/idare/viewmodel"
	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	base "github.com/hayrat/go-template2/backend/pkg/handlers"
)

type CommentHandler struct {
	base.BaseHandler[model.Comment, viewmodel.CommentCreateVM, viewmodel.CommentUpdateVM, viewmodel.CommentListVM, viewmodel.CommentDetailVM]
	CommentService service.ICommentService
}

func NewCommentHandler(s service.ICommentService) CommentHandler {
	h := CommentHandler{
		BaseHandler:    base.NewBaseHandler[model.Comment, viewmodel.CommentCreateVM, viewmodel.CommentUpdateVM, viewmodel.CommentListVM, viewmodel.CommentDetailVM](s),
		CommentService: s,
	}

	return h
}

func (h CommentHandler) GetByID(ctx *app.Ctx) error {
	// Get CommentID from the request (URL parameter in this case)
	CommentID := ctx.ParamsInt64("id") // Assuming the URL is `/Comments/:id`
	if CommentID == 0 {
		return errorsx.BadRequestError("Invalid Comment ID")
	}
	Comment, err := h.CommentService.GetByID(ctx.Context(), CommentID)
	if err != nil {
		return err
	}

	result := viewmodel.CommentMeVM{}.ToViewModel(Comment)

	return ctx.SuccessResponse(result)
}

func (h CommentHandler) MeUpdate(ctx *app.Ctx) error {
	CommentID := ctx.ParamsInt64("id") // Assuming the URL is `/Comments/:id`
	if CommentID == 0 {
		return errorsx.BadRequestError("Invalid Comment ID")
	}
	m, err := h.CommentService.GetByID(ctx.Context(), CommentID)
	if err != nil {
		return err
	}

	var vm viewmodel.CommentMeUpdateVM
	if errs := ctx.BodyParseValidate(&vm); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	m = vm.ToDBModel(m)
	err = h.CommentService.Update(ctx.Context(), m)
	if err != nil {
		return err
	}

	return nil
}

func (h CommentHandler) Create(ctx *app.Ctx) error {
	var comment model.Comment
	if err := ctx.BodyParseValidate(&comment); err != nil {
		return errorsx.ValidationError(err)
	}

	// Token'dan user ID'yi al
	comment.UserId = ctx.GetUserID()

	err := h.CommentService.Create(ctx.Context(), &comment)
	if err != nil {
		return err
	}

	return ctx.SuccessResponse(comment)
}
