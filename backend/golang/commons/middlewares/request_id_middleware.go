package middlewares

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func SetRequestId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := uuid.New().String()
		ctx := context.WithValue(c.Request().Context(), RequestIdKey, requestId)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
