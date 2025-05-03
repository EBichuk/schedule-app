package app

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"schedule-app/config"
	server "schedule-app/internal/app/controller/grpc"
	controller "schedule-app/internal/app/controller/http"
	"schedule-app/internal/app/model"
	"schedule-app/internal/app/service"
	"schedule-app/internal/app/storage"
	"schedule-app/internal/app/storage/pg"
	"schedule-app/internal/pkg/logger"
	"schedule-app/internal/pkg/validators"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	"schedule-app/internal/app/middleware"
)

type App struct {
	c          *controller.Controller
	s          *service.Service
	r          *storage.Repository
	echo       *echo.Echo
	grPCServer *grpc.Server
}

func New() (*App, error) {

	configs := config.GetConfig()

	logger := logger.GetLogger(configs.Logs.Logfile)

	db, err := pg.InitDB(configs.Db)
	if err != nil {
		logger.Error("could not load the database")
	}
	err = model.MigrationSchedule(db)
	if err != nil {
		logger.Error("could not migrate db")
	}

	a := &App{}
	a.r = storage.New(db)
	a.s = service.New(a.r, &configs.MedicationPeriod, logger)
	a.c = controller.New(a.s, logger)

	a.echo = echo.New()
	a.echo.Use(middleware.LoggerEchoReqId)

	a.echo.Validator = &validators.CustomValidator{Validator: validator.New()}

	a.echo.GET("/schedules/:user_id", a.c.GetSchedulesByUser)
	a.echo.GET("/schedule/:user_id/:schedule_id", a.c.GetScheduleById)
	a.echo.GET("/next_takings/:user_id", a.c.NextTaking)
	a.echo.POST("/schedule", a.c.CreateSchedule)

	grPCServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryServerInterceptor))
	server.RegisterServerAPI(grPCServer, a.s, logger)
	a.grPCServer = grPCServer

	return a, nil
}

func (a *App) RunHttp() error {
	err := a.echo.Start(":9000")

	slog.Info("http-server is started")

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (a *App) RunGrpc() error {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	slog.Info("gRPC server is starting", slog.String("adress", l.Addr().String()))

	if err := a.grPCServer.Serve(l); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (a *App) RunApp() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := a.RunGrpc()
		if err != nil {
			log.Fatal("Not start gRPC Server", err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		err := a.RunHttp()
		if err != nil {
			log.Fatal("Not start HTTP Server", err)
		}
	}()
	wg.Wait()
}
