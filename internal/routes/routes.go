package routes

import (
	"echo_practice/internal/controllers"
	"echo_practice/internal/middlewares"

	"github.com/labstack/echo/v4"
)

type Deps struct {
	UserController    *controllers.UserController
	ProfileController *controllers.ProfileController
	ArticleController *controllers.ArticleController
	JWTSecret         string
}

func Register(e *echo.Echo, d Deps) {
	api := e.Group("/api")

	api.POST("/users", d.UserController.Register)
	api.POST("/users/login", d.UserController.Login)

	auth := api.Group("", middlewares.JWTAuth(d.JWTSecret))
	auth.GET("/user", d.UserController.GetCurrentUser)
	auth.PUT("/user", d.UserController.UpdateUser)

	optAuth := api.Group("", middlewares.OptionalJWTAuth(d.JWTSecret))
	optAuth.GET("/profiles/:username", d.ProfileController.GetProfile)

	authProfile := api.Group("", middlewares.JWTAuth(d.JWTSecret))
	authProfile.POST("/profiles/:username/follow", d.ProfileController.FollowUser)
	authProfile.DELETE("/profiles/:username/follow", d.ProfileController.UnfollowUser)

	authArticle := api.Group("", middlewares.JWTAuth(d.JWTSecret))
	authArticle.POST("/articles", d.ArticleController.CreateArticle)
	authArticle.PUT("/articles/:slug", d.ArticleController.UpdateArticle)
	authArticle.DELETE("/articles/:slug", d.ArticleController.DeleteArticle)

	optArticle := api.Group("", middlewares.OptionalJWTAuth(d.JWTSecret))
	optArticle.GET("/articles/:slug", d.ArticleController.GetArticle)
}
