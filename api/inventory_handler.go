package api

import (
	"encoding/json"
	"net/http"

	"github.com/fkrhykal/outbox-cdc/internal/inventory/query"
)

func GetProductsHandler(queryHandler query.GetProductQueryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := queryHandler.GetProducts(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}
