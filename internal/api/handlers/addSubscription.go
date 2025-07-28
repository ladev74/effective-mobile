package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"subscriptions/internal/api"
	"subscriptions/internal/storage/postgresClient"
)

func AddSubscriptionHandler(logger *zap.Logger, pc postgresClient.PostgresClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		subscription := &api.Subscription{}

		err := json.NewDecoder(r.Body).Decode(subscription)
		if err != nil {
			writeResponseWithError(logger, w, http.StatusBadRequest, err.Error())
			logger.Error("AddSubscriptionHandler: cannot decode body", zap.Error(err))
			return
		}

		id, err := pc.SaveSubscription(subscription)
		if err != nil {
			writeResponseWithError(logger, w, http.StatusInternalServerError, err.Error())
			logger.Error("AddSubscriptionHandler:", zap.Error(err))
			return
		}

		writeJSONResponse(logger, w, http.StatusCreated, id)
	}
}
