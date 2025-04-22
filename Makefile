gen:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o /gen/gen.server.go -generate echo-server -package api docs/openapi.yaml
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o /gen/gen.models.go -generate types -package api docs/openapi.yaml