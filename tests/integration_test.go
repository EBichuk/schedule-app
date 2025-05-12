package tests

import (
	"fmt"
	"net/http"
	"schedule-app/internal/app"
	"schedule-app/internal/config"
	"schedule-app/internal/domain/service"
	"schedule-app/internal/infrastructure/persistence"
	server "schedule-app/internal/server/grpc"
	httpserver "schedule-app/internal/server/http"
	"schedule-app/pkg/dbtest"
	grpc_service "schedule-app/pkg/grpc/gen"
	"schedule-app/pkg/logger"
	"schedule-app/pkg/middleware"

	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	Db          *gorm.DB
	grpcClient  grpc_service.UserServiceClient
	httpClient  *http.Client
	baseHTTPURL string
}

func (s *Suite) SetupSuite() {
	rq := s.Require()

	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", "localhost", "scheduleuser", "12345", "postgres", "5433", "disable")

	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	s.Db = db

	go func() {
		a := s.AppInit()
		a.RunApp()
	}()
	time.Sleep(time.Second * 3)

	grpcConn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	rq.NoError(err)

	s.grpcClient = grpc_service.NewUserServiceClient(grpcConn)

	s.httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	s.baseHTTPURL = "http://localhost:9000"
}

func (s *Suite) SetupTest() {
	rq := s.Require()

	err := dbtest.MigrateFromFile(s.Db, "testdata/cleanup.sql")
	rq.NoError(err)
}

func (s *Suite) TearDownSuite() {
	rq := s.Require()

	d, _ := s.Db.DB()
	rq.NoError(d.Close())
}

func (s *Suite) AppInit() *app.App {
	logger := logger.GetLogger("app.log")

	r := persistence.New(s.Db)
	ser := service.New(r, &config.MedPeriodConfig{Period: "3h", Start: "08:00:00", End: "22:00:00"})
	c := httpserver.New(ser, logger)

	grPCServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryServerInterceptor))
	server.RegisterServerAPI(grPCServer, ser, logger)

	a := app.New(c, ser, r, grPCServer)
	return a
}
