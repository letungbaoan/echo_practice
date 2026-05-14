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
		return err
	}

	var req dto.CreateCommentRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	resp, err := h.svc.AddComment(c.Param("slug"), userID, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *CommentController) ListComments(c echo.Context) error {
	resp, err := h.svc.ListComments(c.Param("slug"), middlewares.CurrentUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *CommentController) DeleteComment(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return apperrors.New(http.StatusUnprocessableEntity, "invalid comment id")
	}

	if err := h.svc.DeleteComment(c.Param("slug"), uint(id), userID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{})
}
