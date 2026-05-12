package controllers

import (
	"net/http"
	"strconv"

	"echo_practice/internal/apperrors"
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
		return errorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	var req dto.CreateArticleRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, "invalid json body")
	}
	if err := c.Validate(&req); err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	article, err := h.svc.CreateArticle(userID, req)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, article)
}

func (h *ArticleController) GetArticle(c echo.Context) error {
	slug := c.Param("slug")

	var currentUserID *uint
	if val := c.Get(middlewares.CtxUserID); val != nil {
		if uid, ok := val.(uint); ok {
			currentUserID = &uid
		}
	}

	article, err := h.svc.GetArticle(slug, currentUserID)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, "article not found")
		}
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, article)
}

func (h *ArticleController) UpdateArticle(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return errorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	slug := c.Param("slug")

	var req dto.UpdateArticleRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, "invalid json body")
	}
	if err := c.Validate(&req); err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	article, err := h.svc.UpdateArticle(slug, userID, req)
	if err != nil {
		switch {
		case apperrors.Is(err, apperrors.ErrNotFound):
			return errorResponse(c, http.StatusNotFound, "article not found")
		case apperrors.Is(err, apperrors.ErrForbidden):
			return errorResponse(c, http.StatusForbidden, "forbidden")
		default:
			return errorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, article)
}

func (h *ArticleController) ListArticles(c echo.Context) error {
	filter := services.ListArticlesFilter{
		Tag:               c.QueryParam("tag"),
		AuthorUsername:    c.QueryParam("author"),
		FavoritedUsername: c.QueryParam("favorited"),
		Limit:             parseIntQuery(c, "limit", 0),
		Offset:            parseIntQuery(c, "offset", 0),
	}

	var currentUserID *uint
	if val := c.Get(middlewares.CtxUserID); val != nil {
		if uid, ok := val.(uint); ok {
			currentUserID = &uid
		}
	}

	resp, err := h.svc.ListArticles(filter, currentUserID)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *ArticleController) FeedArticles(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return errorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	limit := parseIntQuery(c, "limit", 0)
	offset := parseIntQuery(c, "offset", 0)

	resp, err := h.svc.FeedArticles(userID, limit, offset)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *ArticleController) FavoriteArticle(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return errorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	resp, err := h.svc.FavoriteArticle(c.Param("slug"), userID)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, "article not found")
		}
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *ArticleController) UnfavoriteArticle(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return errorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	resp, err := h.svc.UnfavoriteArticle(c.Param("slug"), userID)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, "article not found")
		}
		return errorResponse(c, http.StatusInternalServerError, err.Error())
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

func (h *ArticleController) DeleteArticle(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return errorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	slug := c.Param("slug")

	err = h.svc.DeleteArticle(slug, userID)
	if err != nil {
		switch {
		case apperrors.Is(err, apperrors.ErrNotFound):
			return errorResponse(c, http.StatusNotFound, "article not found")
		case apperrors.Is(err, apperrors.ErrForbidden):
			return errorResponse(c, http.StatusForbidden, "forbidden")
		default:
			return errorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{})
}
