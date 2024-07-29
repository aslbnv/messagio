package server

import (
	"encoding/json"
	"net/http"

	"github.com/aslbnv/messagio/internal/types"
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
