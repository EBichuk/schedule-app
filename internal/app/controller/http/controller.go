package controller

import (
	"log/slog"
	"net/http"
	"schedule-app/internal/app/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Service interface {
	CreateSchedule(*model.Schedule) (*model.Schedule, error)
	GetUsersSchedules(uint64) ([]uint64, error)
	GetScheduleByScheduleId(uint64, uint64) (*model.ScheduleTo, error)
	NextTaking(uint64) ([]model.ScheduleTo, error)
}

type Controller struct {
	s Service
}

func New(s Service) *Controller {
	return &Controller{
		s: s,
	}
}

func (c *Controller) GetSchedulesByUser(ctx echo.Context) error {
	userId, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		slog.Info("invalid user id")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	schedulesByUser, err := c.s.GetUsersSchedules(userId)
	if err != nil {
		slog.Info(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, schedulesByUser)
}

func (c *Controller) CreateSchedule(ctx echo.Context) error {
	var schedule model.Schedule

	err := ctx.Bind(&schedule)
	if err != nil {
		slog.Info(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err = ctx.Validate(&schedule)
	if err != nil {
		slog.Info("invalid input data schedule")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input data")
	}

	createdSchedule, err := c.s.CreateSchedule(&schedule)
	if err != nil {
		slog.Info(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, createdSchedule)
}

func (c *Controller) GetScheduleById(ctx echo.Context) error {
	scheduleId, err := strconv.ParseUint(ctx.Param("schedule_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid schudule id")
	}

	userId, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	scheduleById, err := c.s.GetScheduleByScheduleId(scheduleId, userId)
	if err != nil {
		slog.Info(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, scheduleById)
}

func (c *Controller) NextTaking(ctx echo.Context) error {
	userId, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	schedulesByUser, err := c.s.NextTaking(userId)
	if err != nil {
		slog.Info(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, schedulesByUser)
}
