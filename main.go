package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fkrhykal/outbox-cdc/bootstrap"

	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()

	bootstrap.OrderService(mux)

	server := &http.Server{
		Addr:    ":9000",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v\n", err)
	}

	log.Println("server gracefully stopped")
}
