BINARY=./bin/
STD_PATH=cmd/main.go

.PHONY: all user-service

all: user-service

user-service:
	 cd user-service/cmd/ && go build -o ../../bin/user-service

product-service:
	cd porduct-service/cmd/ && go build -o ../../bin/product-service

clean:
	rm -rf $(BINARY)

docker-restart:
	docker rm -f private-go-test-task_postgres_1 private-go-test-task_product-service_1 private-go-test-task_user-service_1
	docker rmi private-go-test-task_product-service private-go-test-task_user-service
	docker-compose up --build -d
