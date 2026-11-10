.PHONY: all
all: lint unit-test integration-test

.PHONY: compile-proto
compile-proto:
	protoc \
		--go_out=. \
		--go-grpc_out=. \
		--proto_path=../../common/proto \
		../../common/proto/metadata.proto

.PHONY: lint
lint:
	golangci-lint run ./... --config ./.golangci.yml

.PHONY: unit-test
unit-test:
	go test -v -cover ./...

.PHONY: integration-test
integration-test:
	go test -v -cover ./... -tags=integration -timeout 3m

.PHONY: coverage-html
coverage-html:
	go test -coverprofile=coverage.out -tags=integration ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: coverage-total
coverage-total:
	go test -coverprofile=coverage.out -tags=integration ./...
	go tool cover -func=coverage.out
	rm coverage.out
