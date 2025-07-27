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
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("AddSubscriptionHandler:", zap.Error(err))
			return
		}

		id, err := pc.SaveSubscription(subscription)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("AddSubscriptionHandler:", zap.Error(err))
			return
		}

		writeResponseWithId(logger, w, http.StatusCreated, id)
	}
}

//curl -X POST http://localhost:8081/add-subscription -H "Content-Type: application/json" \
//-d '{
//"service_name":"yandex plus",
//"price":400,
//"user_id":"123",
//"start_date":"07-2025",
//"end_date":"08-2025"
//}'
