.PHONY: all
all: lint unit-test integration-test

.PHONY: lint
lint:
	golangci-lint run ./... --config ./.golangci.yml

.PHONY: unit-test
unit-test:
	go test -v -race -cover ./...

.PHONY: integration-test
integration-test:
	go test -v -race -cover ./... -tags=integration -timeout 3m

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
