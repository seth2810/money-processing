# Environment
FROM golang:1.16-alpine3.14

WORKDIR /go/src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN apk add build-base

CMD ["go", "test", "-v", "--tags=integration", "./tests"]
