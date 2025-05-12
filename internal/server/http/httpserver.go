package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"schedule-app/internal/domain/entity"
	"strconv"

	"github.com/labstack/echo/v4"
)

type service interface {
	CreateSchedule(context.Context, *entity.Schedule) (*entity.Schedule, error)
	GetUsersSchedules(context.Context, int64) ([]int64, error)
	GetScheduleByScheduleId(context.Context, int64, int64) (*entity.ScheduleWithTime, error)
	NextTaking(context.Context, int64) ([]entity.ScheduleWithTime, error)
}

type HttpServer struct {
	s      service
	logger *slog.Logger
}

func New(s service, logger *slog.Logger) *HttpServer {
	return &HttpServer{
		s:      s,
		logger: logger,
	}
}

func (c *HttpServer) GetSchedulesByUser(ctx echo.Context) error {
	userId, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		c.logger.ErrorContext(ctx.Request().Context(), "invalid user id")
		return ctx.JSON(http.StatusBadRequest, "invalid user id")
	}

	schedulesByUser, err := c.s.GetUsersSchedules(ctx.Request().Context(), userId)
	if err != nil {
		c.logger.ErrorContext(ctx.Request().Context(), err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, schedulesByUser)
}

func (c *HttpServer) CreateSchedule(ctx echo.Context) error {
	var schedule entity.Schedule

	err := ctx.Bind(&schedule)
	if err != nil {
		c.logger.ErrorContext(ctx.Request().Context(), err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err = ctx.Validate(&schedule)
	if err != nil {
		c.logger.ErrorContext(ctx.Request().Context(), "invalid input data schedule")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input data")
	}

	createdSchedule, err := c.s.CreateSchedule(ctx.Request().Context(), &schedule)

	if err != nil {
		c.logger.ErrorContext(ctx.Request().Context(), err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, createdSchedule)
}

func (c *HttpServer) GetScheduleById(ctx echo.Context) error {
	scheduleId, err := strconv.ParseInt(ctx.Param("schedule_id"), 10, 64)
	if err != nil {
		c.logger.ErrorContext(ctx.Request().Context(), "invalid user id")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid schudule id")
	}

	userId, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		c.logger.ErrorContext(ctx.Request().Context(), "invalid user id")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	scheduleById, err := c.s.GetScheduleByScheduleId(ctx.Request().Context(), scheduleId, userId)
	if err != nil {
		c.logger.ErrorContext(ctx.Request().Context(), err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, scheduleById)
}

func (c *HttpServer) NextTaking(ctx echo.Context) error {
	userId, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		c.logger.ErrorContext(ctx.Request().Context(), "invalid user id")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	schedulesByUser, err := c.s.NextTaking(ctx.Request().Context(), userId)
	if err != nil {
		c.logger.ErrorContext(ctx.Request().Context(), err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, schedulesByUser)
}
