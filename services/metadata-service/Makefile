.PHONY: compile-proto lint

compile-proto:
	protoc --go_out=. --go-grpc_out=. --proto_path=../../common/proto ../../common/proto/metadata.proto

lint:
	golangci-lint run ./... --config ./.golangci.yml
