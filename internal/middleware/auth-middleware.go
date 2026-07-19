package middleware

import (
	"net/http"
	"strings"
	"context"
	"github.com/VLKasabiev/simple-wallet/internal/model"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := model.ParseToken(tokenString, secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}


			ctx := context.WithValue(c.Request().Context(), "userID", claims.UserID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}