.PHONY: help
help:
	@echo List of commands:
	@echo   unit-test               - run unit-tests
	@echo   docker-up               - docker compose up
	@echo   docker-down             - docker compose down
	@echo Usage:
	@echo                           make `cmd_name`

.PHONY: unit-test
unit-test:
	go test -cover ./internal/...

.PHONY: graph-gen
graph-gen:
	go run github.com/99designs/gqlgen generate

.PHONY: docker-up
docker-up:
	docker compose up

.PHONY: docker-down
docker-down:
	docker compose down
