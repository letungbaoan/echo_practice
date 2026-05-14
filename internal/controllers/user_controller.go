package controllers

import (
	"net/http"

	"echo_practice/internal/dto"
	"echo_practice/internal/middlewares"
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
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	user, token, err := h.svc.Register(req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, services.ToUserResponse(user, token))
}

func (h *UserController) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	user, token, err := h.svc.Login(req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, services.ToUserResponse(user, token))
}

func (h *UserController) GetCurrentUser(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	user, err := h.svc.GetCurrentUser(userID)
	if err != nil {
		return err
	}

	token, err := h.svc.GenerateToken(user.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, services.ToUserResponse(user, token))
}

func (h *UserController) UpdateUser(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	var req dto.UpdateRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.svc.UpdateUser(userID, req)
	if err != nil {
		return err
	}

	token, err := h.svc.GenerateToken(user.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, services.ToUserResponse(user, token))
}
