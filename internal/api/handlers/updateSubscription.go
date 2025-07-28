package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"subscriptions/internal/api"
	"subscriptions/internal/storage/postgresClient"
)

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
