package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// response uses in writeJSONResponse and writeResponseWithError for structured response to HTTP client.
type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// parseIdParam extracts and validates the id URL parameters from request path.
func parseIdParam(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// writeJSONResponse sets header as application/json and writes response to HTTP client with specified data and status code.
func writeJSONResponse(logger *zap.Logger, w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := response{
		Status: http.StatusText(statusCode),
		Data:   data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Warn("writeJSONResponse: failed to encode response", zap.Error(err))
	}
}

// writeResponseWithError sets header as application/json and writes response about error.
func writeResponseWithError(logger *zap.Logger, w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := response{
		Message: message,
		Status:  http.StatusText(statusCode),
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Warn("writeResponseWithError: cannot send report to caller", zap.Error(err))
	}
}
