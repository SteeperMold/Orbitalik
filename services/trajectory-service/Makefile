.PHONY: dev-update-proto
dev-update-proto:
	cp -r ../../common/proto .

.PHONY: all
all: lint unit-test integration-test

.PHONY: lint
lint:
	cargo fmt --all -- --check
	cargo clippy --all-targets --all-features -- -D warnings

.PHONY: unit-test
unit-test:
	cargo test --features grpc-test-server

.PHONY: integration-test
integration-test:
	cargo test --test '*' --all-features

.PHONY: test
test:
	cargo test --all-features

.PHONY: coverage-html
coverage-html:
	cargo llvm-cov --all-features --html --open

.PHONY: coverage-total
coverage-total:
	cargo llvm-cov --all-features

.PHONY: coverage
coverage: coverage-total coverage-html