CWD := ${shell pwd}

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20

.PHONY: vendor test

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## Run server dev
dev: 
	docker compose  -f docker-compose.development.yml up -d --force-recreate
	make seed
	@go run github.com/cosmtrek/air --build.cmd "go build -o tmp cmd/server/main.go"

## Docker compose up for development
up:
	@echo "Read env from configs/secret.yaml"
	export $(cat configs/secret.yaml | grep -v '#' | xargs)
	@docker compose  -f docker-compose.yml up -d --force-recreate

## Upgrade dependencies
upgrade:
	@go get -u
	@go mod tidy

## Migration database
migration:
	@go run cmd/migration/main.go

## Seed database
seed:
	@go run cmd/seed/main.go

## Mock database args 
mock:
	$(eval ARGS := $(filter-out $@,$(MAKECMDGOALS)))
	@go run cmd/mock/main.go $(ARGS)

## Setup test environment
setup-test:
	@docker compose -f docker-compose.test.yml up -d
	@sleep 5
	
## Down test environment
down-test:
	@docker compose -f docker-compose.test.yml down

## Run all tests
test: #only services package
	@go test -p 1 -v -cover -short ./... 

## Run all tests with coverage
test-cover:  #only services package
	@go test -p 1 -coverprofile=coverage.out ./... 
	@go tool cover -html=coverage.out

# Run docker service
run:
	@docker compose up -d


# Kill port if it's already in use
kill-%:
	@lsof -ti :$* | xargs kill -9