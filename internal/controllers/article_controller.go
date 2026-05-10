package controllers

import (
	"net/http"

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
