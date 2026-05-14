package controllers

import (
	"net/http"

	"echo_practice/internal/services"

	"github.com/labstack/echo/v4"
)

type TagController struct {
	svc *services.TagService
}

func NewTagController(svc *services.TagService) *TagController {
	return &TagController{svc: svc}
}

func (h *TagController) ListTags(c echo.Context) error {
	resp, err := h.svc.ListTags()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}
