package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Err struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"error": message,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
