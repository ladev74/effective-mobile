package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"effmob/internal/api/decoder"
	"effmob/internal/storage/postgresClient"
)

func AddSubscriptionHandler(logger *zap.Logger, pc postgresClient.PostgresClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		subscription, err := decoder.DecodeRequest(logger, r.Body)
		if err != nil {
			writeResponseWithError(logger, w, http.StatusBadRequest, err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("AddSubscriptionHandler:", zap.Error(err))
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
