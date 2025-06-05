package handlers

import (
	"net/http"

	"Weather-Forecast-API/internal/repository"
	"Weather-Forecast-API/internal/utilities"
	"github.com/go-chi/chi/v5"
)

func Confirm(writer http.ResponseWriter, request *http.Request) {
	token := chi.URLParam(request, "token")

	if token == "" {
		utilities.RespondJSON(writer, http.StatusBadRequest, "Invalid token")

		return
	}

	err := repository.ConfirmByToken(token)
	if err != nil {
		if err.Error() == "not found" {
			utilities.RespondJSON(writer, http.StatusNotFound, "Token not found")
		} else {
			utilities.RespondJSON(writer, http.StatusBadRequest, "Failed to get weather: "+err.Error())
		}

		return
	}

	utilities.RespondJSON(writer, http.StatusOK, "Subscription confirmed successfully.")
}
