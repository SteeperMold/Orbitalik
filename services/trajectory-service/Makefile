.PHONY: lint
lint:
	cargo fmt -- --check
	cargo clippy --all-targets --all-features -- -D warnings -W clippy::pedantic -W clippy::nursery

.PHONY: dev-rebuild-proto
dev-update-proto:
	@echo "Copying proto files..."
	cp -r ../../proto .