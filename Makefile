gen:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o /gen/gen.server.go -generate echo-server -package api docs/openapi.yaml
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o /gen/gen.models.go -generate types -package api docs/openapi.yaml

protogen:
	protoc -I proto proto/schedule.proto --go_out=./proto/gen/ --go_opt=paths=source_relative --go-grpc_out=./proto/gen/ --go-grpc_opt=paths=source_relative