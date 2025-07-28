package handlers

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"subscriptions/internal/storage/postgresClient"
)

// GetSubscriptionHandler godoc
// @Summary Get subscription by ID
// @Description Returns subscription details for the given subscription ID.
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} api.Subscription
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Failure 500 {object} response
// @Router /subscriptions/{id} [get]
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
