package controllers

import (
	"net/http"

	"echo_practice/internal/apperrors"
	"echo_practice/internal/dto"
	"echo_practice/internal/services"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	svc *services.UserService
}

func NewUserController(svc *services.UserService) *UserController {
	return &UserController{svc: svc}
}

func (h *UserController) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, "invalid json body")
	}
	if err := c.Validate(&req); err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	user, token, err := h.svc.Register(req)
	if err != nil {
		switch {
		case apperrors.Is(err, apperrors.ErrEmailTaken):
			return errorResponse(c, http.StatusUnprocessableEntity, "email already taken")
		case apperrors.Is(err, apperrors.ErrUsernameTaken):
			return errorResponse(c, http.StatusUnprocessableEntity, "username already taken")
		default:
			return errorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, services.ToUserResponse(user, token))
}

func (h *UserController) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, "invalid json body")
	}
	if err := c.Validate(&req); err != nil {
		return errorResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	user, token, err := h.svc.Login(req)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrInvalidLogin) {
			return errorResponse(c, http.StatusUnauthorized, "invalid email or password")
		}
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, services.ToUserResponse(user, token))
}

func errorResponse(c echo.Context, status int, msg string) error {
	return c.JSON(status, map[string]any{
		"errors": map[string]any{
			"body": []string{msg},
		},
	})
}
