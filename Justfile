set windows-shell := ["cmd.exe", "/c"]

db: 
	docker compose --file docker-compose.yml up --detach --build

dbmigrate: 
  docker exec -i schedule_app_db psql -h 127.0.0.1 -U scheduleuser -d postgres < db/init_db.sql

app: 
	just db
	just dbmigrate

gen:
	oapi-codegen -generate types -package api docs/openapi.yaml > pkg/types/types.gen.go

protogen:
	protoc -I proto pkg/grpc/schedule.proto --go_out=./pkg/grpc/gen/ --go_opt=paths=source_relative --go-grpc_out=./pkg/grpc/gen/ --go-grpc_opt=paths=source_relative

test-infrastructure: 
	docker compose --file tests/docker-compose.yml up --detach --build

test-infrastructure-down:
	docker compose --file tests/docker-compose.yml down --remove-orphans

test-integration: test-unit
	go test -v ./tests -coverprofile tests/cover.out -coverpkg=./...
	
test-unit:
	go test -v ./internal/domain/service -coverprofile tests/cover.out -coverpkg=./...

migrate: 
	docker exec -i schedule_app_db_test psql -h 127.0.0.1 -U scheduleuser -d postgres < db/init_db.sql

test:	
	just test-infrastructure
	just migrate
	just test-integration || (just test-infrastructure-down && exit 1)


build:
  go build -o bin/app ./cmd/main.go

run:
  go run ./cmd/main.go

install-dept:
	go install -u gorm.io/gorm@v1.25.12
	go install gorm.io/driver/postgres@v1.5.11
	go install github.com/labstack/echo/v4@v4.13.3
