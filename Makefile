DOCKER_MIN_VERSION := 19.03.0
DOCKER_COMPOSE_MIN_VERSION := 1.25.0

DOCKER_VERSION := $(shell docker --version | awk '{print $$3}' | sed 's/,//')
DOCKER_COMPOSE_VERSION := $(shell docker-compose --version | awk '{print $$3}' | sed 's/,//')

.PHONY: all tests

all: tests

check-version:
	@echo Docker version is $(DOCKER_VERSION)
	@echo Docker Compose version is $(DOCKER_COMPOSE_VERSION)

run:
	docker-compose up --build -d
stop:
	docker-compose down
tests:
	cd user-service/internal/service && go test -v -cover
	cd user-service/internal/repository && go test -v -cover
	cd recommendation-service/internal/service && go test -v -cover
	cd recommendation-service/internal/repository && go test -v -cover
	cd analytics-service/internal/repository && go test -v -cover
	cd product-service/internal/service && go test -v -cover
	cd product-service/internal/repository && go test -v -cover
