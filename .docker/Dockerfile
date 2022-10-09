# Собираем в гошке
FROM golang:1.19-alpine3.16 as build-env

WORKDIR /go/src

COPY . .

RUN CGO_ENABLED=0 go build -o /opt/service/serve ./cmd/serve

# Release
FROM alpine:3.16

WORKDIR /usr/service

COPY --from=build-env /opt/service/serve /usr/service/serve
COPY --from=build-env /go/src/migrations /usr/service/migrations

CMD ["/usr/service/serve"]
