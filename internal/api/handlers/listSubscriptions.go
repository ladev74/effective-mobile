package handlers

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"effmob/internal/storage/postgresClient"
)

func ListSubscriptionsHandler(logger *zap.Logger, pc postgresClient.PostgresClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		subscriptions, err := pc.ListSubscriptions()
		if err != nil {
			switch {
			case errors.Is(err, postgresClient.ErrSubscriptionNotFound):
				writeResponseWithError(logger, w, http.StatusNotFound, err.Error())
			default:
				writeResponseWithError(logger, w, http.StatusInternalServerError, err.Error())
			}

			logger.Error("ListSubscriptionsHandler:", zap.Error(err))
			return
		}

		writeJSONResponse(logger, w, http.StatusOK, subscriptions)
	}
}
