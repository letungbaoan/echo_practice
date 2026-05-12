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

type CommentController struct {
	svc *services.CommentService
}

func NewCommentController(svc *services.CommentService) *CommentController {
	return &CommentController{svc: svc}
}

func (h *CommentController) AddComment(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return errorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	slug := c.Param("slug")

	var req dto.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, "invalid json body")
	}
	if err := c.Validate(&req); err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	resp, err := h.svc.AddComment(slug, userID, req)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, "article not found")
		}
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *CommentController) ListComments(c echo.Context) error {
	slug := c.Param("slug")

	var currentUserID *uint
	if val := c.Get(middlewares.CtxUserID); val != nil {
		if uid, ok := val.(uint); ok {
			currentUserID = &uid
		}
	}

	resp, err := h.svc.ListComments(slug, currentUserID)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, "article not found")
		}
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *CommentController) DeleteComment(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return errorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	slug := c.Param("slug")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, "invalid comment id")
	}

	err = h.svc.DeleteComment(slug, uint(id), userID)
	if err != nil {
		switch {
		case apperrors.Is(err, apperrors.ErrNotFound):
			return errorResponse(c, http.StatusNotFound, "not found")
		case apperrors.Is(err, apperrors.ErrForbidden):
			return errorResponse(c, http.StatusForbidden, "forbidden")
		default:
			return errorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, map[string]any{})
}
