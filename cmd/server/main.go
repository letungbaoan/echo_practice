package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"echo_practice/internal/config"
	"echo_practice/internal/controllers"
	"echo_practice/internal/database"
	"echo_practice/internal/repositories"
	"echo_practice/internal/routes"
	"echo_practice/internal/services"
	"echo_practice/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("db migrate: %v", err)
	}
	log.Println("database connected & migrated")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo, cfg.JWTSecret)
	userCtrl := controllers.NewUserController(userSvc)

	e := echo.New()
	e.Validator = utils.NewValidator()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []any{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
			}
			if v.Error != nil {
				attrs = append(attrs, slog.String("error", v.Error.Error()))
				logger.Error("request", attrs...)
			} else {
				logger.Info("request", attrs...)
			}
			return nil
		},
	}))
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	routes.Register(e, routes.Deps{
		UserController: userCtrl,
	})

	log.Fatal(e.Start(":" + cfg.Port))
}
