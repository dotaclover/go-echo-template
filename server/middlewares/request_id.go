package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const RequestIDKey = "request_id"

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			c.Set(RequestIDKey, requestID)
			c.Response().Header().Set(echo.HeaderXRequestID, requestID)
			return next(c)
		}
	}
}
