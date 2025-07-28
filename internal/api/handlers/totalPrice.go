package handlers

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"subscriptions/internal/api"
	"subscriptions/internal/storage/postgresClient"
)

func GetTotalPriceHandler(logger *zap.Logger, pc postgresClient.PostgresClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, serviceName, startDate, endDate := getQuery(r)

		startPeriod, err := time.Parse("01-2006", startDate)
		if err != nil {
			writeResponseWithError(logger, w, http.StatusBadRequest, "invalid 'from' date format. Use MM-YYYY")
			return
		}

		endPeriod, err := time.Parse("01-2006", endDate)
		if err != nil {
			writeResponseWithError(logger, w, http.StatusBadRequest, "invalid 'to' date format. Use MM-YYYY")
			return
		}

		subscriptions, err := pc.ListFilteredSubscriptions(userID, serviceName)
		if err != nil {
			logger.Error("GetTotalCostHandler: failed to load subscriptions", zap.Error(err))
			writeResponseWithError(logger, w, http.StatusInternalServerError, "cannot load subscriptions")
			return
		}

		total := processSubscriptions(startPeriod, endPeriod, subscriptions)

		writeJSONResponse(logger, w, http.StatusOK, total)
	}
}

func getQuery(r *http.Request) (string, string, string, string) {
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	return userID, serviceName, startDate, endDate
}

func processSubscriptions(startPeriod time.Time, endPeriod time.Time, subscriptions []*api.Subscription) int {
	total := 0
	for _, subscription := range subscriptions {
		subscriptionStart, err := time.Parse("01-2006", subscription.StartDate)
		if err != nil {
			continue
		}

		subscriptionEnd := endPeriod
		if subscription.EndDate != "" {
			subscriptionEnd, err = time.Parse("01-2006", subscription.EndDate)
			if err != nil {
				continue
			}
		}

		var intersectStart time.Time

		if subscriptionStart.After(startPeriod) {
			intersectStart = subscriptionStart
		} else {
			intersectStart = startPeriod
		}

		var intersectEnd time.Time

		if subscriptionEnd.Before(endPeriod) {
			intersectEnd = subscriptionEnd
		} else {
			intersectEnd = endPeriod
		}

		if !intersectStart.After(intersectEnd) {
			years := intersectEnd.Year() - intersectStart.Year()
			months := int(intersectEnd.Month()) - int(intersectStart.Month())
			months = years*12 + months + 1
			total += months * subscription.Price
		}
	}

	return total
}
