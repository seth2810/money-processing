package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/seth2810/money-processing/internal/server"
)

func main() {
	ctx, cancelFn := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
	)

	defer cancelFn()

	if err := server.Start(ctx, "0.0.0.0:8080"); err != nil {
		log.Fatalln("Error:", err)
	}
}
