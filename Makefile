.PHONY: run-dev
run-dev:
	@docker compose up --build
	@docker compose down

.PHONY: down
down:
	@docker compose down --remove-orphans

.PHONY: flush
flush:
	@docker compose down --volumes --remove-orphans

.PHONY: prepare-env
prepare-env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
	else \
		echo "==> Skipping, .env already exists"; \
	fi
	@cd services; \
	for service in $(SERVICES); do \
		if [ ! -f $$service/.env ]; then \
			cp $$service/.env.example $$service/.env; \
		else \
			echo "==> Skipping $$service (.env already exists)"; \
		fi \
	done
