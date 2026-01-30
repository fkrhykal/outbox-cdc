package bootstrap

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fkrhykal/outbox-cdc/api"
	"github.com/fkrhykal/outbox-cdc/data"
	"github.com/fkrhykal/outbox-cdc/db"
	"github.com/fkrhykal/outbox-cdc/internal/order/service"
)

func OrderService(mux *http.ServeMux) {
	pg, err := sql.Open("postgres", "user=pguser password=pgpw dbname=pgdb port=5432 host=localhost sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	orderRepository := db.NewPgOrderRepository(pg)
	outboxRepository := db.NewPgOutboxRepository(pg)

	txManager := data.NewSqlTxManager(pg)

	orderService := service.NewOrderService(
		txManager,
		orderRepository,
		outboxRepository,
	)

	mux.Handle("POST /orders", api.PlaceOrderHandler(orderService))
}
