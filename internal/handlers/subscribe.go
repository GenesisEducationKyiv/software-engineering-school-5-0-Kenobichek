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

func Subscribe(writer http.ResponseWriter, request *http.Request) {
	channelValue := request.FormValue("email")
	city := request.FormValue("city")
	frequency := request.FormValue("frequency")

	if channelValue == "" || city == "" || frequency == "" {
		utilities.RespondJSON(writer, http.StatusBadRequest, "Invalid input")

		return
	}

	channelType := request.FormValue("channelType")
	if channelType == "" {
		channelType = "email"
	}

	if !utilities.IsValidChannel(channelType) {
		utilities.RespondJSON(writer, http.StatusBadRequest, "Unsupported channelType")

		return
	}

	frequencyMinutes, err := utilities.ConvertFrequency(frequency)
	if err != nil {
		utilities.RespondJSON(writer, http.StatusBadRequest, err.Error())

		return
	}

	template, err := repository.GetTemplateByName("confirm")
	if err != nil {
		utilities.RespondJSON(writer, http.StatusInternalServerError, "Failed to load confirmation template")

		return
	}

	token := uuid.NewString()

	sub := &models.Subscription{
		ChannelType:      channelType,
		ChannelValue:     channelValue,
		City:             city,
		FrequencyMinutes: frequencyMinutes,
		Token:            token,
	}

	if err := repository.CreateSubscription(sub); err != nil {
		utilities.RespondJSON(writer, http.StatusConflict, "Already subscribed or DB error")

		return
	}

	message := strings.ReplaceAll(template.Message, "{{ confirm_token }}", token)
	subject := template.Subject

	emailNotifier := notifier.EmailNotifier{}
	_ = emailNotifier.Send(channelValue, message, subject)

	utilities.RespondJSON(writer, http.StatusOK, "Subscription successful. Confirmation sent.")
}
