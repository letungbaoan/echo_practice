package middlewares

import (
	"fmt"
	"net/http"

	"echo_practice/internal/apperrors"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	status, messages := mapError(err)
	_ = c.JSON(status, map[string]any{
		"errors": map[string]any{
			"body": messages,
		},
	})
}

func mapError(err error) (int, []string) {
	if appErr, ok := apperrors.As(err); ok {
		return appErr.Status, appErr.Messages
	}

	var he *echo.HTTPError
	if asHTTPError(err, &he) {
		return he.Code, []string{fmt.Sprint(he.Message)}
	}

	return http.StatusInternalServerError, []string{"internal server error"}
}

func asHTTPError(err error, target **echo.HTTPError) bool {
	if e, ok := err.(*echo.HTTPError); ok {
		*target = e
		return true
	}
	return false
}
