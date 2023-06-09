MOCKERY_VERSION = 2.23.1
SERVICE_VERSION = 0.0.1
SERVICE_IMAGE = class-booking-service

install-tools:
	go install github.com/vektra/mockery/v2@v${MOCKERY_VERSION}

compile:
	go build ./...

clean down:
	docker-compose -f ./scripts/docker/docker-compose-test.yml down -v

test-run-int:
	INTEGRATION=1 go test -count=1 ./...

test-unit:
	go test -count=1 ./...

test-int-up:
	docker-compose -f ./scripts/docker/docker-compose-test.yml up -d --remove-orphans
	POSTGRES_PASSWORD=class_booking POSTGRES_USER=class_booking bash ./scripts/docker/wait-for-postgres.sh localhost

build:
	go build -o .bin/class-booking-service -v ./app/services/booking

# clear local databases before running test
test-int: down test-int-up test-run-int clean

#test-int runs unit and integration
test: test-int

#Local
start-local:
	docker-compose -f ./scripts/docker/docker-compose.yml up -d --remove-orphans postgres
	POSTGRES_PASSWORD=class_booking POSTGRES_USER=class_booking bash ./scripts/docker/wait-for-postgres.sh localhost

stop-local:
	docker-compose -f ./scripts/docker/docker-compose.yml down

run: start-local
	go run -v ./app/services/booking

build-docker:
	docker build -f scripts/docker/Dockerfile -t $(SERVICE_IMAGE) .

