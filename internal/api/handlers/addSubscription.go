package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"effmob/internal/api/decoder"
	"effmob/internal/storage/postgresClient"
)

func NewAddSubscription(logger *zap.Logger, pc postgresClient.PostgresClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		subscription, err := decoder.DecodeRequest(logger, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("NewAddSubscription:", zap.Error(err))
			return
		}

		id, err := pc.SaveSubscription(subscription)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("NewAddSubscription:", zap.Error(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// TODO: добавить перенос строки, вынести в отдельную функцию (создать структуру статус?)
		//err := json.NewEncoder(w).Encode(resp)

		err = json.NewEncoder(w).Encode(struct {
			id int `json:"id"`
		}{
			id: id,
		})
		if err != nil {
			logger.Warn("NewAddSubscription: cannot send report to caller", zap.Error(err))
		}
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
