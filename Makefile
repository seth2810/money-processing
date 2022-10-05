DOCKER_IMG="banners_rotation:develop"

build:
	go build -v -o ./bin/serve ./cmd/serve

build-img:
	docker build -t "money-processing:develop" -f .docker/Dockerfile .

test:
	go test -race -count 100 ./internal/...
