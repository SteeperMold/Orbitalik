.PHONY: compile-proto lint

compile-proto:
	protoc --go_out=. --go-grpc_out=. --proto_path=../../proto ../../proto/tle.proto

lint:
	golangci-lint run ./... --config ./.golangci.yml
