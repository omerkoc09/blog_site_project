package handlers

import (
	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/common/service"
	"github.com/hayrat/go-template2/backend/idare/viewmodel"
	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	base "github.com/hayrat/go-template2/backend/pkg/handlers"
)

type LikeHandler struct {
	base.BaseHandler[model.Like, viewmodel.LikeCreateVM, viewmodel.LikeUpdateVM, viewmodel.LikeListVM, viewmodel.LikeDetailVM]
	LikeService         service.ILikeService
	NotificationService service.INotificationService
	PostService         service.IPostService
}

func NewLikeHandler(s service.ILikeService, ns service.INotificationService, ps service.IPostService) LikeHandler {
	h := LikeHandler{
		BaseHandler:         base.NewBaseHandler[model.Like, viewmodel.LikeCreateVM, viewmodel.LikeUpdateVM, viewmodel.LikeListVM, viewmodel.LikeDetailVM](s),
		LikeService:         s,
		NotificationService: ns,
		PostService:         ps,
	}

	return h
}

func (h LikeHandler) GetByID(ctx *app.Ctx) error {
	// Get LikeID from the request (URL parameter in this case)
	LikeID := ctx.ParamsInt64("id") // Assuming the URL is `/Likes/:id`
	if LikeID == 0 {
		return errorsx.BadRequestError("Invalid Like ID")
	}
	Like, err := h.LikeService.GetByID(ctx.Context(), LikeID)
	if err != nil {
		return err
	}

	result := viewmodel.LikeMeVM{}.ToViewModel(Like)

	return ctx.SuccessResponse(result)
}

func (h LikeHandler) MeUpdate(ctx *app.Ctx) error {
	LikeID := ctx.ParamsInt64("id") // Assuming the URL is `/Likes/:id`
	if LikeID == 0 {
		return errorsx.BadRequestError("Invalid Like ID")
	}
	m, err := h.LikeService.GetByID(ctx.Context(), LikeID)
	if err != nil {
		return err
	}

	var vm viewmodel.LikeMeUpdateVM
	if errs := ctx.BodyParseValidate(&vm); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	m = vm.ToDBModel(m)
	err = h.LikeService.Update(ctx.Context(), m)
	if err != nil {
		return err
	}

	return nil
}

func (h LikeHandler) Create(ctx *app.Ctx) error {
	postID := ctx.ParamsInt64("id")
	if postID == 0 {
		return errorsx.BadRequestError("Invalid post ID")
	}
	userID := ctx.GetUserID()
	err := h.LikeService.ToggleLike(ctx.Context(), userID, postID)
	if err != nil {
		return err
	}
	// Bildirim: post sahibine
	post, err := h.PostService.GetByID(ctx.Context(), postID)
	if err == nil && post.UserId != userID {
		notification := model.Notification{
			SenderID:     userID,
			TargetUserID: post.UserId,
			Type:         model.NotificationTypeLike,
			PostID:       &postID,
			IsRead:       false,
		}
		_ = h.NotificationService.Create(ctx.Context(), &notification)
	}
	return ctx.SuccessResponse(nil)
}
