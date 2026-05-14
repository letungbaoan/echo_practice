package middlewares

import (
	"strings"

	"echo_practice/internal/apperrors"
	"echo_practice/internal/utils"

	"github.com/labstack/echo/v4"
)

const CtxUserID = "userID"

func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr, ok := extractToken(c)
			if !ok {
				return apperrors.ErrUnauthorized
			}

			claims, err := utils.ParseToken(tokenStr, jwtSecret)
			if err != nil || claims == nil {
				return apperrors.ErrUnauthorized
			}

			c.Set(CtxUserID, claims.UserID)
			return next(c)
		}
	}
}

func OptionalJWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if tokenStr, ok := extractToken(c); ok {
				if claims, err := utils.ParseToken(tokenStr, jwtSecret); err == nil && claims != nil {
					c.Set(CtxUserID, claims.UserID)
				}
			}
			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (uint, error) {
	val := c.Get(CtxUserID)
	if val == nil {
		return 0, apperrors.ErrUnauthorized
	}
	userID, ok := val.(uint)
	if !ok {
		return 0, apperrors.ErrUnauthorized
	}
	return userID, nil
}

func CurrentUserID(c echo.Context) *uint {
	val := c.Get(CtxUserID)
	if val == nil {
		return nil
	}
	uid, ok := val.(uint)
	if !ok {
		return nil
	}
	return &uid
}

func extractToken(c echo.Context) (string, bool) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Token" {
		return "", false
	}
	return parts[1], true
}
