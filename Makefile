APP_NAME=companies
CUR_DIR=$(shell pwd)
SRC=$(CUR_DIR)/cmd/companies
BINARY_NAME=$(CUR_DIR)/bin/$(APP_NAME)

.PHONY: generate fmt test test/int build run clean mod/tidy build-container run-container diagrams

all: fmt generate test build clean mod/tidy

generate:
	$(info ************ GENERATE MOCKS ************)
	go generate -v ./...

fmt:
	$(info ************ RUN FROMATING ************)
	go fmt ./...

test:
	$(info ************ RUN UNIT TESTS ************)
	go test -count=1 -v ./...

test/int:
	$(info ************ RUN UNIT AND INTEGRATION TESTS ************)
	go test -count=1 -tags integration -v ./...

build:
	$(info ************ BUILD ************)
	CGO_ENABLED=0 go build -o $(BINARY_NAME) -v $(SRC)

run:
	$(info ************ RUN ************)
	$(BINARY_NAME)

lint:
	$(info ************ Linter ************)
	golangci-lint run ./... -v

clean:
	$(info ************ CLEAN ************)
	go clean
	rm -f $(BINARY_NAME)

mod/tidy:
	$(info ************ MOD TIDY ************)
	go mod tidy

build/docker:
	$(info ************ BUILD CONTAINER ************)
	docker build -t $(APP_NAME) .
run/docker:
	$(info ************ RUN CONTAINER ************)
	docker run --rm --env DB_HOST=host.docker.internal -p 8080:8080 --name $(APP_NAME) $(APP_NAME)

migrate/up:
	$(info ************ MIGRATE UP ************)
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up
migrate/down:
	$(info ************ MIGRATE DOWN ************)
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose down 1

up:
	$(info ************ DOCKER-COMPOSE UP ************)
	docker-compose up -d
up/postgres:
	$(info ************ UP POSTGRES IN DOCKER-COMPOSE ************)
	docker-compose up -d postgres
down:
	$(info ************ DOCKER-COMPOSE DOWN ************)
	docker-compose down

kafka/topic/create:
	$(info ************ Kafka topic create ************)
	docker exec -it companies-kafka-1 /usr/bin/kafka-topics --create --topic companies-mutations --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka/topic/describe:
	$(info ************ Kafka topic describe ************)
	docker exec -it companies-kafka-1 /usr/bin/kafka-topics --describe --topic companies-mutations --bootstrap-server localhost:9092
kafka/topic/consume:
	$(info ************ Kafka console consumer ************)
	docker exec -it companies-kafka-1 /usr/bin/kafka-console-consumer --topic companies-mutations --from-beginning --bootstrap-server localhost:9092

test/int/docker-compose: up kafka/topic/create test/int down

run/local: up kafka/topic/create build run

company/create:
	$(info ************ Create company ************)
	curl --location 'http://localhost:8080/api/v1/secured/companies' \
    --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6InVrZXNoLm11cnVnYW4iLCJlbWFpbCI6InVrZXNoQGdvLmNvbSIsImV4cCI6MTcxMzM4NzYwMH0.Dbcz0odhXAbdjM5HprynZ4eSv-OCBZqhymCZOC-MKiM' \
    --header 'Content-Type: application/json' \
    --data '{ \
        "id": "abc8c242-00ed-40a6-82df-ea0d3afd0867", \
        "name": "XM67", \
        "employees_amount": 123, \
        "registered": true, \
        "type": "Corporations" \
    }' -w "\n\n"

company/patch:
	$(info ************ Patch company ************)
	curl --location --request PATCH 'http://localhost:8080/api/v1/secured/companies/abc8c242-00ed-40a6-82df-ea0d3afd0867' \
    --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6InVrZXNoLm11cnVnYW4iLCJlbWFpbCI6InVrZXNoQGdvLmNvbSIsImV4cCI6MTcxMzM4NzYwMH0.Dbcz0odhXAbdjM5HprynZ4eSv-OCBZqhymCZOC-MKiM' \
    --header 'Content-Type: application/json' \
    --data '{ \
        "name": "XM67", \
        "description": "description67", \
        "employees_amount": 66, \
        "registered": false, \
        "type": "Sole Proprietorship" \
    }' -w "\n\n"

company/delete:
	$(info ************ Delete company ************)
	curl --location --request DELETE 'http://localhost:8080/api/v1/secured/companies/abc8c242-00ed-40a6-82df-ea0d3afd0867' \
    --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6InVrZXNoLm11cnVnYW4iLCJlbWFpbCI6InVrZXNoQGdvLmNvbSIsImV4cCI6MTcxMzM4NzYwMH0.Dbcz0odhXAbdjM5HprynZ4eSv-OCBZqhymCZOC-MKiM'  -w "\n\n"

company/get:
	$(info ************ Get company ************)
	curl --location 'http://localhost:8080/api/v1/companies/abc8c242-00ed-40a6-82df-ea0d3afd0867' -w "\n\n"

company/livecycle: company/create company/get company/patch company/get company/delete

diagrams:
	$(info ************ GENERATE DIAGRAMS ************)
	@if command -v dot >/dev/null 2>&1; then \
		find docs/diagrams -name "*.dot" -exec sh -c 'dot -Tpng $$1 -o $${1%.dot}.png' _ {} \; ; \
		echo "Диаграммы успешно сгенерированы"; \
	else \
		echo "Graphviz не установлен. Установите его для генерации диаграмм:"; \
		echo "  macOS: brew install graphviz"; \
		echo "  Ubuntu/Debian: sudo apt-get install graphviz"; \
		echo "  CentOS/RHEL: sudo yum install graphviz"; \
	fi
