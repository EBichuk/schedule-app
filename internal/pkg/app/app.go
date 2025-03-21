package app

import (
	"log"
	"schedule-app/config"
	"schedule-app/internal/app/controller"
	"schedule-app/internal/app/model"
	"schedule-app/internal/app/service"
	"schedule-app/internal/app/storage"
	"schedule-app/internal/app/storage/pg"
	"schedule-app/internal/pkg/validators"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type App struct {
	c    *controller.Controller
	s    *service.Service
	r    *storage.Repository
	echo *echo.Echo
}

func New() (*App, error) {

	configs := config.GetConfig()

	db, err := pg.InitDB(configs.Db)
	if err != nil {
		log.Fatal("could not load the database")
	}
	err = model.MigrationSchedule(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	a := &App{}
	a.r = storage.New(db)
	a.s = service.New(a.r, &configs.MedicationPeriod)
	a.c = controller.New(a.s)

	a.echo = echo.New()
	a.echo.Validator = &validators.CustomValidator{Validator: validator.New()}

	a.echo.GET("/schedules/:user_id", a.c.GetSchedulesByUser)
	a.echo.GET("/schedule/:user_id/:schedule_id", a.c.GetScheduleById)
	a.echo.GET("/next_takings/:user_id", a.c.NextTaking)
	a.echo.POST("/schedule", a.c.CreateSchedule)

	return a, nil
}

func (a *App) Run() error {
	err := a.echo.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
