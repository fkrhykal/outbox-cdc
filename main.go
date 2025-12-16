package main

import (
	"database/sql"
	"log"
	"net/http"

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

	http.ListenAndServe(":9000", mux)
}
