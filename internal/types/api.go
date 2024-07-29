package types

import (
	"net/http"
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error
