package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/seth2810/money-processing/internal/config"
	"github.com/seth2810/money-processing/internal/routes"
	"github.com/seth2810/money-processing/internal/server"
	"github.com/seth2810/money-processing/internal/storage"
)

func initDB(ctx context.Context, cfg *config.DatabaseConfig) (*sql.DB, error) {
	dbDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := goose.OpenDBWithDriver("postgres", dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return nil, fmt.Errorf("failed to migrate DB: %w", err)
	}

	return db, nil
}

func serve(ctx context.Context) error {
	cfg, err := config.Read()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	db, err := initDB(ctx, &cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}

	defer db.Close()

	s := storage.New(db)

	srv := &http.Server{
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      routes.NewRouter(s),
		Addr:         net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
	}

	return server.Start(ctx, srv)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func main() {
	ctx, cancelFn := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
	)

	defer cancelFn()

	checkErr(serve(ctx))
}
