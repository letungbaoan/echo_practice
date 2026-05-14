package controllers

import (
	"net/http"
	"strconv"

	"echo_practice/internal/dto"
	"echo_practice/internal/middlewares"
	"echo_practice/internal/services"

	"github.com/labstack/echo/v4"
)

type ArticleController struct {
	svc *services.ArticleService
}

func NewArticleController(svc *services.ArticleService) *ArticleController {
	return &ArticleController{svc: svc}
}

func (h *ArticleController) CreateArticle(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateArticleRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	article, err := h.svc.CreateArticle(userID, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, article)
}

func (h *ArticleController) GetArticle(c echo.Context) error {
	article, err := h.svc.GetArticle(c.Param("slug"), middlewares.CurrentUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, article)
}

func (h *ArticleController) UpdateArticle(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	var req dto.UpdateArticleRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	article, err := h.svc.UpdateArticle(c.Param("slug"), userID, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, article)
}

func (h *ArticleController) DeleteArticle(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	if err := h.svc.DeleteArticle(c.Param("slug"), userID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{})
}

func (h *ArticleController) ListArticles(c echo.Context) error {
	filter := services.ListArticlesFilter{
		Tag:               c.QueryParam("tag"),
		AuthorUsername:    c.QueryParam("author"),
		FavoritedUsername: c.QueryParam("favorited"),
		Limit:             parseIntQuery(c, "limit", 0),
		Offset:            parseIntQuery(c, "offset", 0),
	}

	resp, err := h.svc.ListArticles(filter, middlewares.CurrentUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *ArticleController) FeedArticles(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	resp, err := h.svc.FeedArticles(userID, parseIntQuery(c, "limit", 0), parseIntQuery(c, "offset", 0))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *ArticleController) FavoriteArticle(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	resp, err := h.svc.FavoriteArticle(c.Param("slug"), userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *ArticleController) UnfavoriteArticle(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	resp, err := h.svc.UnfavoriteArticle(c.Param("slug"), userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func parseIntQuery(c echo.Context, key string, def int) int {
	s := c.QueryParam(key)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
