package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"subscriptions/internal/api"
	"subscriptions/internal/storage/postgresClient"
)

// UpdateSubscriptionHandler godoc
// @Summary Update a subscription by ID
// @Description Updates subscription data for the given subscription ID.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body api.Subscription true "Subscription data to update"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /subscriptions/{id} [put]
func UpdateSubscriptionHandler(logger *zap.Logger, pc postgresClient.PostgresClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIdParam(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("UpdateSubscriptionHandler: cannot get id from URL", zap.Error(err))
			return
		}

		subscription := &api.Subscription{}

		err = json.NewDecoder(r.Body).Decode(subscription)
		if err != nil {
			writeResponseWithError(logger, w, http.StatusBadRequest, err.Error())
			logger.Error("UpdateSubscriptionHandler: cannot decode body", zap.Error(err))
			return
		}

		err = pc.UpdateSubscription(id, subscription)
		if err != nil {
			writeResponseWithError(logger, w, http.StatusInternalServerError, err.Error())
			logger.Error("UpdateSubscriptionHandler:", zap.Error(err))
			return
		}

		writeJSONResponse(logger, w, http.StatusOK, nil)
	}
}
