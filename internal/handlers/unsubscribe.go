package handlers

import (
	"Weather-Forecast-API/internal/notifier"
	"Weather-Forecast-API/internal/repository"
	"Weather-Forecast-API/internal/utilities"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if token == "" {
		utilities.RespondJSON(w, http.StatusBadRequest, "Invalid input")
		return
	}

	template, err := repository.GetTemplateByName("unsubscribe")
	if err != nil {
		utilities.RespondJSON(w, http.StatusInternalServerError, "Failed to load unsubscribe template")
		return
	}

	subscription, err := repository.GetSubscriptionByToken(token)
	if err != nil {
		utilities.RespondJSON(w, http.StatusConflict, err.Error())
		return
	}

	err = repository.UnsubscribeByToken(token)
	if err != nil {
		if err.Error() == "not found" {
			utilities.RespondJSON(w, http.StatusNotFound, "Token not found")
		} else {
			utilities.RespondJSON(w, http.StatusBadRequest, "Failed to get weather: "+err.Error())
		}
		return
	}

	message := strings.ReplaceAll(template.Message, "{{ city }}", subscription.City)
	subject := template.Subject

	notifier := notifier.EmailNotifier{}
	_ = notifier.Send(subscription.ChannelValue, message, subject)

	utilities.RespondJSON(w, http.StatusOK, "You have been unsubscribed successfully.")
}
