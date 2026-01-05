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
