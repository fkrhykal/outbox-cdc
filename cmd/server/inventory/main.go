package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fkrhykal/outbox-cdc/api"
	"github.com/fkrhykal/outbox-cdc/api/middleware"
	"github.com/fkrhykal/outbox-cdc/db"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/service"
	_ "github.com/lib/pq"
)

func main() {

	mux := http.NewServeMux()

	pg, err := sql.Open("postgres", "user=pguser password=pgpw dbname=inventory_db port=5432 host=localhost sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := pg.Ping(); err != nil {
		log.Fatal(err)
	}

	productRepository := db.NewPgProductRepository(pg)

	productService := service.NewProductService(productRepository)

	mux.Handle("GET /products", api.GetProductsHandler(productService))

	server := &http.Server{
		Addr:    ":8000",
		Handler: middleware.LogMiddleware(mux),
	}

	log.Println("Listening on http://localhost:8000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s: %v\n", server.Addr, err)
	}

	if err := pg.Close(); err != nil {
		log.Fatalf("failed to close database connection: %v", err)
	}
}
