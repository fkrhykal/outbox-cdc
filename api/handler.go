package api

import (
	"encoding/json"
	"net/http"

	"github.com/fkrhykal/outbox-cdc/internal/command"
)

func PlaceOrderHandler(h command.PlaceOrderHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cmd := new(command.PlaceOrder)

		if err := json.NewDecoder(r.Body).Decode(cmd); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		res, err := h.PlaceOrder(r.Context(), cmd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})
}
