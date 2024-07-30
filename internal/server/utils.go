package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aslbnv/messagio/internal/types"
	"github.com/gorilla/mux"
)

func makeHandler(fn types.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, types.APIError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func getID(r *http.Request) (int, error ) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}

	return id, nil
}
