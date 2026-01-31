package api

import (
	"encoding/json"
	"net/http"

	command "github.com/fkrhykal/outbox-cdc/internal/order/comand"
	"github.com/fkrhykal/outbox-cdc/internal/validation"
)

func PlaceOrderHandler(commandHandler command.PlaceOrderHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd := new(command.PlaceOrder)

		if err := json.NewDecoder(r.Body).Decode(cmd); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}

		res, err := commandHandler.PlaceOrder(r.Context(), cmd)

		if err == nil {
			w.Header().Add("content-type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(res); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
			return
		}

		if validationErrors, ok := err.(validation.Errors); ok {
			w.Header().Add("content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(validationErrors); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
