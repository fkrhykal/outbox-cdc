package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fkrhykal/outbox-cdc/api"
	"github.com/fkrhykal/outbox-cdc/data"
	"github.com/fkrhykal/outbox-cdc/db"
	"github.com/fkrhykal/outbox-cdc/internal/outbox"
	"github.com/fkrhykal/outbox-cdc/internal/service"
	_ "github.com/lib/pq"
)

func main() {

	pg, err := sql.Open("postgres", "user=pguser password=pgpw dbname=pgdb port=5432 host=localhost sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	orderRepository := db.NewPgOrderRepository(pg)
	outboxRepository := db.NewPgOutboxRepository(pg)

	txManager := data.NewSqlTxManager(pg)
	publisher := outbox.NewOutboxEventPublisher(outboxRepository)

	orderService := service.NewOrderService(
		txManager,
		orderRepository,
		publisher,
	)

	mux := http.NewServeMux()

	mux.Handle("POST /orders", api.PlaceOrderHandler(orderService))

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
