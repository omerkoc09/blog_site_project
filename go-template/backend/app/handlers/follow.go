package handlers

import (
	"github.com/hayrat/go-template2/backend/app/viewmodel"
	"github.com/hayrat/go-template2/backend/common/service"
	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
)

type FollowHandler struct {
	FollowService service.IFollowService
}

func NewFollowHandler(s service.IFollowService) FollowHandler {
	h := FollowHandler{
		FollowService: s,
	}

	return h
}

func (h FollowHandler) Follow(ctx *app.Ctx) error {
	followingID := ctx.ParamsInt64("id")
	if followingID == 0 {
		return errorsx.BadRequestError("Invalid following ID")
	}

	followerID := ctx.GetUserID()

	follow, err := h.FollowService.Follow(ctx.Context(), followerID, followingID)
	if err != nil {
		return err
	}

	// Create a new viewmodel and populate it from the follow model
	result := viewmodel.FollowDetailVM{}.ToViewModel(follow)
	return ctx.SuccessResponse(result)
}

func (h FollowHandler) Unfollow(ctx *app.Ctx) error {
	followingID := ctx.ParamsInt64("id")
	if followingID == 0 {
		return errorsx.BadRequestError("Invalid following ID")
	}

	followerID := ctx.GetUserID()

	err := h.FollowService.Unfollow(ctx.Context(), followerID, followingID)
	if err != nil {
		return err
	}

	// Return a simple success response
	return ctx.SuccessResponse(map[string]bool{"success": true})
}

// Add new methods to handle getting followers and following

func (h FollowHandler) GetFollowers(ctx *app.Ctx) error {
	userID := ctx.ParamsInt64("id")
	if userID == 0 {
		return errorsx.BadRequestError("Invalid user ID")
	}

	users, err := h.FollowService.GetFollowers(ctx.Context(), userID)
	if err != nil {
		return err
	}

	// Convert users to the appropriate view model
	userInfos := make([]viewmodel.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = viewmodel.UserToUserInfo(user)
	}

	result := viewmodel.FollowersListVM{
		Count:     len(userInfos),
		Followers: userInfos,
	}

	return ctx.SuccessResponse(result)
}

func (h FollowHandler) GetFollowing(ctx *app.Ctx) error {
	userID := ctx.ParamsInt64("id")
	if userID == 0 {
		return errorsx.BadRequestError("Invalid user ID")
	}

	users, err := h.FollowService.GetFollowing(ctx.Context(), userID)
	if err != nil {
		return err
	}

	// Convert users to the appropriate view model
	userInfos := make([]viewmodel.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = viewmodel.UserToUserInfo(user)
	}

	result := viewmodel.FollowingListVM{
		Count:     len(userInfos),
		Following: userInfos,
	}

	return ctx.SuccessResponse(result)
}
