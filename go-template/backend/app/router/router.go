package router

import (
	"time"

	"github.com/hayrat/go-template2/backend/app/handlers"
	"github.com/hayrat/go-template2/backend/common/service"
	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/router"
)

type AppRouter struct {
}

func NewAppRouter() *AppRouter {
	return &AppRouter{}
}

func (AppRouter) RegisterRoutes(app *app.App) {

	api := app.FiberApp.Group("/api")
	secretKey := app.Cfg.Server.JwtSecret

	authService := service.NewAuthService(app.DB, app.Cfg.Server.JwtSecret, app.Cfg.Server.JwtAccessTokenExpireMinute*time.Hour, app.Cfg.Server.JwtRefreshTokenExpireHour*time.Hour)
	userService := service.NewUserService(app.DB)
	postService := service.NewPostService(app.DB)
	tokenService := service.NewTokenService(secretKey)
	commentService := service.NewCommentService(app.DB)
	likeService := service.NewLikeService(app.DB)
	followService := service.NewFollowService(app.DB)
	savedService := service.NewSavedService(app.DB)
	topicService := service.NewTopicService(app.DB)

	authHandler := handlers.NewAuthHandler(authService, userService, *tokenService)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)
	likeHandler := handlers.NewLikeHandler(likeService)
	followHandler := handlers.NewFollowHandler(followService)
	savedHandler := handlers.NewSavedHandler(savedService)
	topicHandler := handlers.NewTopicHandler(topicService)

	router.Post(api, "/auth/login", authHandler.Login)
	router.Post(api, "/auth/refresh", authHandler.RefreshToken)
	router.Post(api, "/auth/logout", authHandler.Logout)

	router.Post(api, "/auth/forgot-password", authHandler.TokenForForgotPassword)
	router.Put(api, "/auth/reset-password", authHandler.ResetPassword)

	router.Post(api, "/user", userHandler.Create)
	router.Get(api, "/post", postHandler.Query)
	router.Get(api, "/post/:id", postHandler.GetByID)
	router.Get(api, "/user", userHandler.Query)
	router.Get(api, "/user_guest/:id", userHandler.GetByID)
	router.Get(api, "/followers/:id", followHandler.GetFollowers)
	router.Get(api, "/followings/:id", followHandler.GetFollowing)

	api.Use(router.JWTMiddleware(app))

	router.Get(api, "/user/me", userHandler.Me)
	router.Put(api, "/user/me", userHandler.MeUpdate)
	router.Get(api, "/user/:id", userHandler.GetByID)

	router.Post(api, "/post", postHandler.Create)
	router.Put(api, "/post/:id", postHandler.MeUpdateWithImage)
	router.Delete(api, "/post/:id", postHandler.Delete)
	router.Post(api, "/post/image", postHandler.CreatePostWithImage)

	router.Post(api, "/comment", commentHandler.Create)
	router.Get(api, "/comment/:id", commentHandler.GetByID)
	router.Put(api, "/comment/:id", commentHandler.MeUpdate)
	router.Delete(api, "/comment/:id", commentHandler.Delete)

	router.Post(api, "/like/:id", likeHandler.Create)
	router.Put(api, "/like/:id", likeHandler.MeUpdate)
	router.Delete(api, "/like/:id", likeHandler.Delete)

	// Register follow routes
	router.Post(api, "/follow/:id", followHandler.Follow)
	router.Delete(api, "/follow/:id", followHandler.Unfollow)
	router.Get(api, "/following/:id", followHandler.GetFollowing)

	// Register saved routes
	router.Post(api, "/saved/:id", savedHandler.ToggleSaved)
	router.Get(api, "/saved", savedHandler.GetSavedPosts)
	router.Get(api, "/saved/:id", savedHandler.IsSaved)

	// Topic endpoints
	router.Post(api, "/topic", topicHandler.Create)
	router.Get(api, "/topic", topicHandler.Query)
	router.Get(api, "/topic/:id", topicHandler.GetByID)
	router.Put(api, "/topic/:id", topicHandler.Update)
	router.Delete(api, "/topic/:id", topicHandler.Delete)

	router.Get(api, "/post/:id/topics", postHandler.GetTopicsByPostID)

}
