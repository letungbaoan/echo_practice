package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"echo_practice/internal/utils"
	"github.com/labstack/echo/v4"
)

const (
	CtxUserID = "userID"
	AuthError = "auth error"
)

func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"errors": map[string][]string{
						"body": {"unauthorized"},
					},
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Token" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"errors": map[string][]string{
						"body": {"unauthorized"},
					},
				})
			}

			tokenStr := parts[1]
			claims, err := utils.ParseToken(tokenStr, jwtSecret)
			if err != nil || claims == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"errors": map[string][]string{
						"body": {"unauthorized"},
					},
				})
			}

			c.Set(CtxUserID, claims.UserID)
			return next(c)
		}
	}
}

func OptionalJWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Token" {
					tokenStr := parts[1]
					if claims, err := utils.ParseToken(tokenStr, jwtSecret); err == nil && claims != nil {
						c.Set(CtxUserID, claims.UserID)
					}
				}
			}
			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (uint, error) {
	val := c.Get(CtxUserID)
	if val == nil {
		return 0, errors.New("user not authenticated")
	}
	userID, ok := val.(uint)
	if !ok {
		return 0, errors.New("invalid user id in context")
	}
	return userID, nil
}
