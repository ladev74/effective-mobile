package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Id      int    `json:"id,omitempty"`
}

func writeResponseWithId(logger *zap.Logger, w http.ResponseWriter, statusCode, id int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := response{
		Status: http.StatusText(statusCode),
		Id:     id,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Warn("AddSubscriptionHandler: cannot send report to caller", zap.Error(err))
	}
}

func writeResponseWithError(logger *zap.Logger, w http.ResponseWriter, statusCode int, message string, id int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := response{
		Message: message,
		Id:      id,
		Status:  http.StatusText(statusCode),
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Warn("writeResponseWithError: cannot send report to caller", zap.Error(err))
	}
}
