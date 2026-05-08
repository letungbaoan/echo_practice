package routes

import (
	"echo_practice/internal/controllers"

	"github.com/labstack/echo/v4"
)

type Deps struct {
	UserController *controllers.UserController
}

func Register(e *echo.Echo, d Deps) {
	api := e.Group("/api")

	api.POST("/users", d.UserController.Register)
	api.POST("/users/login", d.UserController.Login)
}
