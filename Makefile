BIN="./bin/money_processing"
DOCKER_IMG="money_processing:develop"

prepare:
	cp -n .env.example config.env || true

lint:
	golangci-lint run ./...

test:
	go test -race -count 100 ./internal/...

build:
	go build -v -o $(BIN) ./cmd/serve

start: prepare build
	$(BIN)

build-img:
	docker build -t $(DOCKER_IMG) -f .docker/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

up:
	docker-compose up -d --build

down:
	docker-compose down
