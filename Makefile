
.PHONY: all tests

all: tests

tests:
	cd user-service/internal/service && go test -v -cover
	cd user-service/internal/repository && go test -v -cover
	cd recommendation-service/internal/service && go test -v -cover
	cd recommendation-service/internal/repository && go test -v -cover
	cd analytics-service/internal/repository && go test -v -cover
	cd product-service/internal/service && go test -v -cover
	cd product-service/internal/repository && go test -v -cover