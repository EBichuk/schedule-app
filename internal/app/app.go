package app

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"schedule-app/internal/config"
	"schedule-app/internal/domain/service"
	"schedule-app/internal/infrastructure/integration"
	"schedule-app/internal/infrastructure/persistence"
	server "schedule-app/internal/server/grpc"
	httpserver "schedule-app/internal/server/http"
	"schedule-app/pkg/logger"
	"schedule-app/pkg/validators"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	"schedule-app/pkg/middleware"
)

type App struct {
	c          *httpserver.HttpServer
	s          *service.Service
	r          *persistence.Repository
	grPCServer *grpc.Server
}

func New(c *httpserver.HttpServer, s *service.Service, r *persistence.Repository, grPCServer *grpc.Server) *App {
	return &App{
		c:          c,
		s:          s,
		r:          r,
		grPCServer: grPCServer,
	}
}

func Init() (*App, error) {

	configs := config.GetConfig()

	logger := logger.GetLogger(configs.Logs.Logfile)

	db, err := integration.InitDB(configs.Db)
	if err != nil {
		logger.Error("could not load the database")
	}

	a := &App{}
	a.r = persistence.New(db)
	a.s = service.New(a.r, &configs.MedicationPeriod)
	a.c = httpserver.New(a.s, logger)

	grPCServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryServerInterceptor))
	server.RegisterServerAPI(grPCServer, a.s, logger)
	a.grPCServer = grPCServer

	return a, nil
}

func (a *App) RunHttp() error {
	e := echo.New()
	e.Use(middleware.LoggerEchoReqId)
	e.Validator = &validators.CustomValidator{Validator: validator.New()}

	e.GET("/schedules/:user_id", a.c.GetSchedulesByUser)
	e.GET("/schedule/:user_id/:schedule_id", a.c.GetScheduleById)
	e.GET("/next_takings/:user_id", a.c.NextTaking)
	e.POST("/schedule", a.c.CreateSchedule)

	err := e.Start(":9000")

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
