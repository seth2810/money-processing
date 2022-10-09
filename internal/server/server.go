package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, srv *http.Server) error {
	errCh := make(chan error, 1)

	log.Printf("server is starting on %s...", srv.Addr)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	var err error

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errCh:
	}

	if !errors.Is(err, context.Canceled) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	log.Println("server is stopping...")

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}
