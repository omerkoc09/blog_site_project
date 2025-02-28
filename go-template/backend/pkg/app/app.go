package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/uptrace/bun"
	"go.uber.org/zap"

	"github.com/hayrat/go-template2/backend/pkg/config"
	"github.com/hayrat/go-template2/backend/pkg/database"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	"github.com/hayrat/go-template2/backend/pkg/log"
	"github.com/hayrat/go-template2/backend/pkg/viewmodel"
)

type IRouter interface {
	RegisterRoutes(app *App)
}

type App struct {
	FiberApp *fiber.App
	DB       *bun.DB
	Cfg      *config.Config
	Ctx      context.Context
}

func (a *App) Stack() {
	panic("unimplemented")
}

// app.go dosyasında
func (a *App) Static(prefix string, root string) {
	// Fiber.Static metodunu kullanarak statik dosya servisi oluşturun
	a.FiberApp.Static(prefix, root)
}

func New(router IRouter, Version, BuildTime string) *App {
	cfg, err := config.Setup()
	if err != nil {
		panic(err)
	}

	fiberApp := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
		ErrorHandler: errorsx.ErrorHandler,
		BodyLimit:    20 * 1024 * 1024,
		//ReadBufferSize: fiber.DefaultReadBufferSize * 2, // Request Header Fields Too Large hatası için
	})

	fiberApp.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:5174",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Content-Type,Authorization,Origin",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Content-Type",
	}))
	fiberApp.Use(logger.New())
	fiberApp.Use(requestid.New(requestid.Config{
		Header: fiber.HeaderXRequestID,
	}))

	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	fiberApp.Get("/api/version", func(c *fiber.Ctx) error {
		return c.SendString("version: " + Version + " - buildtime: " + BuildTime)
	})
	fiberApp.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return err
		}

		if len(c.Response().Body()) == 0 {
			return c.JSON(&viewmodel.ResponseModel{})
		}
		return nil
	})

	db := database.New(cfg.Database)

	app := &App{
		FiberApp: fiberApp,
		DB:       db,
		Cfg:      cfg,
		Ctx:      context.Background(),
	}

	router.RegisterRoutes(app)

	return app
}

func (a *App) Start() {
	l := log.GetLogger("")
	l.SetOptions(zap.AddCallerSkip(-2))
	l.Info("http server başlatılıyor...")
	go func() {
		err := a.FiberApp.Listen(fmt.Sprintf(":%v", 3000))
		if err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	l.Info("Gracefully shutting down...")
	err := a.FiberApp.Shutdown()
	if err != nil {
		l.Error("FiberApp shutdown", zap.Error(err))
	}
	err = a.DB.Close()
	if err != nil {
		l.Error("DB close", zap.Error(err))
	}
	l.Info("Elhamdülillah")
}

func (a *App) Start2() {
	l := log.GetLogger("")
	l.SetOptions(zap.AddCallerSkip(-2))
	l.Info("http server başlatılıyor...")
	go func() {
		err := a.FiberApp.Listen(fmt.Sprintf(":%v", 3001))
		if err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	l.Info("Gracefully shutting down...")
	err := a.FiberApp.Shutdown()
	if err != nil {
		l.Error("FiberApp shutdown", zap.Error(err))
	}
	err = a.DB.Close()
	if err != nil {
		l.Error("DB close", zap.Error(err))
	}
	l.Info("Elhamdülillah")
}
