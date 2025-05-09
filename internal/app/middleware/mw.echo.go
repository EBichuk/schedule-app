package middleware

import (
	"fmt"
	"log/slog"
	"schedule-app/internal/pkg/contexts"
	"time"

	"github.com/labstack/echo/v4"
)

func LoggerEchoReqId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := contexts.WithRequestID(c.Request().Context())
		fmt.Println(ctx)
		c.SetRequest(c.Request().WithContext(ctx))

		slog.InfoContext(ctx,
			"request",
			slog.String("method", c.Request().Method),
			slog.String("path", c.Request().URL.Path),
			slog.String("remote_addr", c.Request().RemoteAddr),
			slog.String("user_agent", c.Request().UserAgent()),
		)
		t1 := time.Now()

		next(c)

		slog.InfoContext(ctx,
			"request completed",
			slog.Int("status", c.Response().Status),
			slog.Int64("bytes", c.Response().Size),
			slog.String("duration", time.Since(t1).String()))
		return nil
	}
}
