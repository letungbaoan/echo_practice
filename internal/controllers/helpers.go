package controllers

import (
	"echo_practice/internal/apperrors"

	"github.com/labstack/echo/v4"
)

func bindAndValidate(c echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return apperrors.ErrInvalidBody
	}
	if err := c.Validate(req); err != nil {
		return apperrors.New(422, err.Error())
	}
	return nil
}
