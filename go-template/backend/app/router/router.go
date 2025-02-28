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

	authHandler := handlers.NewAuthHandler(authService, userService, *tokenService)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)
	likeHandler := handlers.NewLikeHandler(likeService)

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

	api.Use(router.JWTMiddleware(app))

	router.Get(api, "/user/me", userHandler.Me)
	router.Put(api, "/user/me", userHandler.MeUpdate)
	router.Get(api, "/user/:id", userHandler.GetByID)

	router.Post(api, "/post", postHandler.Create)
	router.Put(api, "/post/:id", postHandler.MeUpdateWithImage)
	router.Delete(api, "/post/:id", postHandler.Delete)
	router.Post(api, "/post/:id/image", postHandler.UploadPostImage)
	router.Post(api, "/post/image", postHandler.CreatePostWithImage)

	router.Post(api, "/comment", commentHandler.Create)
	router.Get(api, "/comment/:id", commentHandler.GetByID)
	router.Put(api, "/comment/:id", commentHandler.MeUpdate)
	router.Delete(api, "/comment/:id", commentHandler.Delete)

	router.Post(api, "/like/:id", likeHandler.Create)
	router.Put(api, "/like/:id", likeHandler.MeUpdate)
	router.Delete(api, "/like/:id", likeHandler.Delete)

}
