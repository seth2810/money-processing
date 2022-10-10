BIN="./bin/money_processing"
DOCKER_IMG="seth2810/money_processing:develop"

prepare:
	cp -n .env.example config.env || true

lint:
	golangci-lint run ./...

build:
	go build -v -o $(BIN) ./cmd/serve

start: prepare build
	$(BIN)

up:
	docker-compose up -d

down:
	docker-compose down

integration-tests:
	set -e ;\
	docker-compose -f docker-compose.test.yml up -d --build --scale tests=0;\
	test_status_code=0 ;\
	docker-compose -f docker-compose.test.yml run tests || test_status_code=$$? ;\
	docker-compose -f docker-compose.test.yml down ;\
	exit $$test_status_code ;

integration-tests-cleanup:
	docker-compose -f docker-compose.test.yml down \
		--rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
	docker-compose rm -f
