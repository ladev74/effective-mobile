package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"effmob/internal/storage/postgresClient"
)

func DeleteSubscriptionHandler(logger *zap.Logger, pc postgresClient.PostgresClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIdParam(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("DeleteSubscriptionHandler: cannot get id from URL", zap.Error(err))
			return
		}

		err = pc.DeleteSubscription(id)
		if err != nil {
			switch {
			case errors.Is(err, postgresClient.ErrSubscriptionNotFound):
				writeResponseWithError(logger, w, http.StatusNotFound, err.Error(), id)
			default:
				writeResponseWithError(logger, w, http.StatusInternalServerError, err.Error(), id)
			}

			logger.Error("DeleteSubscriptionHandler:", zap.Error(err))
			return

		}

		writeResponseWithId(logger, w, http.StatusOK, id)
	}
}

func parseIdParam(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}
