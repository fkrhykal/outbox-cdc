package bootstrap

import (
	"database/sql"
	"log"

	"github.com/fkrhykal/outbox-cdc/data"
	"github.com/fkrhykal/outbox-cdc/db"
	"github.com/fkrhykal/outbox-cdc/internal/order/service"
	"github.com/fkrhykal/outbox-cdc/internal/outbox"
)

func BootstrapOrderService() *service.OrderService[data.SqlExecutor] {
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
	return orderService
}
