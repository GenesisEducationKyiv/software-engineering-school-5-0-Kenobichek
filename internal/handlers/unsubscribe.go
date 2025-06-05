package handlers

import (
	"net/http"
	"strings"

	"Weather-Forecast-API/internal/notifier"
	"Weather-Forecast-API/internal/repository"
	"Weather-Forecast-API/internal/utilities"
	"github.com/go-chi/chi/v5"
)

func Unsubscribe(writer http.ResponseWriter, request *http.Request) {
	token := chi.URLParam(request, "token")

	if token == "" {
		utilities.RespondJSON(writer, http.StatusBadRequest, "Invalid input")

		return
	}

	template, err := repository.GetTemplateByName("unsubscribe")
	if err != nil {
		utilities.RespondJSON(writer, http.StatusInternalServerError, "Failed to load unsubscribe template")

		return
	}

	subscription, err := repository.GetSubscriptionByToken(token)
	if err != nil {
		utilities.RespondJSON(writer, http.StatusConflict, err.Error())

		return
	}

	err = repository.UnsubscribeByToken(token)
	if err != nil {
		if err.Error() == "not found" {
			utilities.RespondJSON(writer, http.StatusNotFound, "Token not found")
		} else {
			utilities.RespondJSON(writer, http.StatusBadRequest, "Failed to get weather: "+err.Error())
		}

		return
	}

	message := strings.ReplaceAll(template.Message, "{{ city }}", subscription.City)
	subject := template.Subject

	emailNotifier := notifier.EmailNotifier{}
	_ = emailNotifier.Send(subscription.ChannelValue, message, subject)

	utilities.RespondJSON(writer, http.StatusOK, "You have been unsubscribed successfully.")
}
