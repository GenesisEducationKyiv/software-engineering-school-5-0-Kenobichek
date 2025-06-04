package handlers

import (
	"net/http"
	"strings"

	"Weather-Forecast-API/internal/models"
	"Weather-Forecast-API/internal/notifier"
	"Weather-Forecast-API/internal/repository"
	"Weather-Forecast-API/internal/utilities"

	"github.com/google/uuid"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	channel_value := r.FormValue("email")
	city := r.FormValue("city")
	frequency := r.FormValue("frequency")

	if channel_value == "" || city == "" || frequency == "" {
		utilities.RespondJSON(w, http.StatusBadRequest, "Invalid input")
		return
	}

	channel_type := r.FormValue("channel_type")
	if channel_type == "" {
		channel_type = "email"
	}

	if !utilities.IsValidChannel(channel_type) {
		utilities.RespondJSON(w, http.StatusBadRequest, "Unsupported channel_type")
		return
	}

	frequency_minutes, err := utilities.ConvertFrequency(frequency)
	if err != nil {
		utilities.RespondJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	template, err := repository.GetTemplateByName("confirm")
	if err != nil {
		utilities.RespondJSON(w, http.StatusInternalServerError, "Failed to load confirmation template")
		return
	}

	token := uuid.NewString()

	sub := &models.Subscription{
		ChannelType:      channel_type,
		ChannelValue:     channel_value,
		City:             city,
		FrequencyMinutes: frequency_minutes,
		Token:            token,
	}

	if err := repository.CreateSubscription(sub); err != nil {
		utilities.RespondJSON(w, http.StatusConflict, "Already subscribed or DB error")
		return
	}

	message := strings.ReplaceAll(template.Message, "{{ confirm_token }}", token)
	subject := template.Subject

	notifier := notifier.EmailNotifier{}
	_ = notifier.Send(channel_value, message, subject)

	utilities.RespondJSON(w, http.StatusOK, "Subscription successful. Confirmation sent.")
}
