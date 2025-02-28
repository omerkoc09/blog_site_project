package router

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"go.uber.org/zap"

	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/log"
	"github.com/hayrat/go-template2/backend/pkg/viewmodel"
)

func Get(r fiber.Router, path string, h func(ctx *app.Ctx) error) {
	r.Get(path, ctxWrap(h))
}

func Post(r fiber.Router, path string, h func(ctx *app.Ctx) error) {
	r.Post(path, ctxWrap(h))
}

func Put(r fiber.Router, path string, h func(ctx *app.Ctx) error) {
	r.Put(path, ctxWrap(h))
}

func Delete(r fiber.Router, path string, h func(ctx *app.Ctx) error) {
	r.Delete(path, ctxWrap(h))
}

func ctxWrap(h func(ctx *app.Ctx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid := c.Get("requestid", "")
		logger := log.GetLogger(rid, zap.String("HTTP_METHOD", c.Method()), zap.String("HTTP_PATH", c.Path()))
		c.Locals("logger", &logger)
		return h(&app.Ctx{Ctx: c})
	}
}

func JWTErrorHandler(ctx *fiber.Ctx, err error) error {
	ctx.Status(fiber.StatusUnauthorized)

	if err.Error() == "Missing or malformed JWT" {
		return ctx.JSON(viewmodel.ResponseModel{ErrorMessage: "token malformed"})
	}
	return ctx.JSON(viewmodel.ResponseModel{ErrorMessage: "token expired"})
}

func JWTMiddleware(app *app.App) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(app.Cfg.Server.JwtSecret),
		ErrorHandler: JWTErrorHandler,
	})
}
