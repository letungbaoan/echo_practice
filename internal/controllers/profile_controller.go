package controllers

import (
	"net/http"

	"echo_practice/internal/apperrors"
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
	username := c.Param("username")

	var currentUserID *uint
	if val := c.Get(middlewares.CtxUserID); val != nil {
		if uid, ok := val.(uint); ok {
			currentUserID = &uid
		}
	}

	profile, err := h.svc.GetProfile(username, currentUserID)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, "profile not found")
		}
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *ProfileController) FollowUser(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return errorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	username := c.Param("username")

	profile, err := h.svc.FollowUser(userID, username)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, "profile not found")
		}
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *ProfileController) UnfollowUser(c echo.Context) error {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return errorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	username := c.Param("username")

	profile, err := h.svc.UnfollowUser(userID, username)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNotFound) {
			return errorResponse(c, http.StatusNotFound, "profile not found")
		}
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, profile)
}
