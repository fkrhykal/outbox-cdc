package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fkrhykal/outbox-cdc/api"
	"github.com/fkrhykal/outbox-cdc/api/middleware"
	"github.com/fkrhykal/outbox-cdc/data"
	"github.com/fkrhykal/outbox-cdc/db"
	"github.com/fkrhykal/outbox-cdc/internal/order/service"
	orderValidation "github.com/fkrhykal/outbox-cdc/internal/order/validation"
	"github.com/fkrhykal/outbox-cdc/internal/validation"
	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()

	pg, err := sql.Open("postgres", "user=pguser password=pgpw dbname=order_db port=5433 host=localhost sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := pg.Ping(); err != nil {
		log.Fatal(err)
	}

	validator := validation.NewValidatorRegistry()
	validator.Register(validation.Coerce(orderValidation.ValidatePlaceOrderCommand()))

	orderRepository := db.NewPgOrderRepository(pg)
	outboxRepository := db.NewPgOutboxRepository(pg)

	txManager := data.NewSqlTxManager(pg)

	orderService := service.NewOrderService(
		validator,
		txManager,
		orderRepository,
		outboxRepository,
	)

	mux.Handle("POST /orders", api.PlaceOrderHandler(orderService))

	server := &http.Server{
		Addr:    ":9000",
		Handler: middleware.LogMiddleware(mux),
	}

	log.Println("Listening on http://localhost:9000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s: %v\n", server.Addr, err)
	}

	if err := pg.Close(); err != nil {
		log.Fatalf("failed to close database connection: %v", err)
	}
}
