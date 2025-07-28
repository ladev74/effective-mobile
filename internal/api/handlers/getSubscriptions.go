package handlers

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"subscriptions/internal/storage/postgresClient"
)

// GetSubscriptionHandler returns an HTTP handler for getting subscription by its id.
func GetSubscriptionHandler(logger *zap.Logger, pc postgresClient.PostgresClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIdParam(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("GetSubscriptionsHandler: cannot get id from URL", zap.Error(err))
			return
		}

		subscription, err := pc.GetSubscription(id)
		if err != nil {
			switch {
			case errors.Is(err, postgresClient.ErrSubscriptionNotFound):
				writeResponseWithError(logger, w, http.StatusNotFound, err.Error())
			default:
				writeResponseWithError(logger, w, http.StatusInternalServerError, err.Error())
			}

			logger.Error("GetSubscriptionsHandler:", zap.Error(err))
			return
		}

		writeJSONResponse(logger, w, http.StatusOK, subscription)
	}
}
