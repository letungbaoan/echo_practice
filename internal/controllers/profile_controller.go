package controllers

import (
	"net/http"

	"echo_practice/internal/middlewares"
	"echo_practice/internal/services"

	"github.com/labstack/echo/v4"
)

type ProfileController struct {
	svc *services.ProfileService
}

func NewProfileController(svc *services.ProfileService) *ProfileController {
	return &ProfileController{svc: svc}
}

func (h *ProfileController) GetProfile(c echo.Context) error {
	profile, err := h.svc.GetProfile(c.Param("username"), middlewares.CurrentUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, profile)
}

func (h *ProfileController) FollowUser(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	profile, err := h.svc.FollowUser(userID, c.Param("username"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, profile)
}

func (h *ProfileController) UnfollowUser(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return err
	}

	profile, err := h.svc.UnfollowUser(userID, c.Param("username"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, profile)
}
