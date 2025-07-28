package handlers

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"subscriptions/internal/storage/postgresClient"
)

// DeleteSubscriptionHandler returns an HTTP handler for deleting a subscription by its id.
func DeleteSubscriptionHandler(logger *zap.Logger, pc postgresClient.PostgresClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIdParam(r)
		if err != nil {
			writeResponseWithError(logger, w, http.StatusBadRequest, err.Error())
			logger.Error("DeleteSubscriptionHandler: cannot get id from URL", zap.Error(err))
			return
		}

		err = pc.DeleteSubscription(id)
		if err != nil {
			switch {
			case errors.Is(err, postgresClient.ErrSubscriptionNotFound):
				writeResponseWithError(logger, w, http.StatusNotFound, err.Error())
			default:
				writeResponseWithError(logger, w, http.StatusInternalServerError, err.Error())
			}

			logger.Error("DeleteSubscriptionHandler:", zap.Error(err))
			return
		}

		writeJSONResponse(logger, w, http.StatusOK, nil)
	}
}
