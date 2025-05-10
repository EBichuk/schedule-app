gen:
	oapi-codegen -generate types -package api docs/openapi.yaml > pkg/types/types.gen.go

protogen:
	protoc -I proto proto/schedule.proto --go_out=./proto/gen/ --go_opt=paths=source_relative --go-grpc_out=./proto/gen/ --go-grpc_opt=paths=source_relative