package handlers

import (
	"Weather-Forecast-API/internal/repository"
	"Weather-Forecast-API/internal/utilities"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Confirm(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if token == "" {
		utilities.RespondJSON(w, http.StatusBadRequest, "Invalid token")
		return
	}

	err := repository.ConfirmByToken(token)
	if err != nil {
		if err.Error() == "not found" {
			utilities.RespondJSON(w, http.StatusNotFound, "Token not found")
		} else {
			utilities.RespondJSON(w, http.StatusBadRequest, "Failed to get weather: "+err.Error())
		}
		return
	}

	utilities.RespondJSON(w, http.StatusOK, "Subscription confirmed successfully.")
}
