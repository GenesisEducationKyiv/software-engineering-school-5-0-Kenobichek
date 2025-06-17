package subscribe

import (
	"Weather-Forecast-API/internal/models"
	"Weather-Forecast-API/internal/repository"
	"Weather-Forecast-API/internal/response"
	"Weather-Forecast-API/internal/services/notification"
	"Weather-Forecast-API/internal/services/subscription"
	"strings"

	"github.com/google/uuid"
	"net/http"
)

type SubscribeHandler struct {
	subService   subscription.SubscriptionService
	notifService notification.NotificationService
}

func NewSubscribeHandler(subService subscription.SubscriptionService,
		notifService notification.NotificationService) *SubscribeHandler {
	return &SubscribeHandler{
		subService:   subService,
		notifService: notifService,
	}
}

func (h *SubscribeHandler) Subscribe(writer http.ResponseWriter, request *http.Request) {
	input, err := ParseAndValidateSubscribeInput(request)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	frequencyMinutes, err := ConvertFrequency(input.Frequency)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	template, err := repository.GetTemplateByName("confirm")
	if err != nil {
		response.RespondJSON(writer, http.StatusInternalServerError, "Failed to load confirmation template")
		return
	}

	token := uuid.NewString()

	sub := &models.Subscription{
		ChannelType:      input.ChannelType,
		ChannelValue:     input.ChannelValue,
		City:             input.City,
		FrequencyMinutes: frequencyMinutes,
		Token:            token,
	}

	if err := h.subService.Subscribe(sub); err != nil {
		response.RespondJSON(writer, http.StatusConflict, "Already subscribed or DB error")
		return
	}

	message := strings.ReplaceAll(template.Message, "{{ confirm_token }}", token)

	err = h.notifService.SendMessage(input.ChannelType, input.ChannelValue, message, template.Subject)
	if err != nil {
		response.RespondJSON(writer, http.StatusInternalServerError, "Failed to send message. Error: "+err.Error())
		return
	}

	response.RespondJSON(writer, http.StatusOK, "Subscription successful. Confirmation sent.")
}

func (h *SubscribeHandler) Unsubscribe(writer http.ResponseWriter, request *http.Request) {
	input, err := ParseAndValidateTokenInput(request)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	template, err := repository.GetTemplateByName("unsubscribe")
	if err != nil {
		response.RespondJSON(writer, http.StatusInternalServerError, "Failed to load unsubscribe template")
		return
	}

	sub, err := repository.GetSubscriptionByToken(input.Token)
	if err != nil {
		response.RespondJSON(writer, http.StatusConflict, err.Error())
		return
	}

	if err := h.subService.Unsubscribe(sub); err != nil {
		if err.Error() == "not found" {
			response.RespondJSON(writer, http.StatusNotFound, "Token not found")
		} else {
			response.RespondJSON(writer, http.StatusBadRequest, "Failed to confirm subscription: "+err.Error())
		}
		return
	}

	message := strings.ReplaceAll(template.Message, "{{ city }}", sub.City)

	err = h.notifService.SendMessage(sub.ChannelType, sub.ChannelValue, message, template.Subject)
	if err != nil {
		response.RespondJSON(writer, http.StatusInternalServerError, "Failed to send message. Error: "+err.Error())
		return
	}

	response.RespondJSON(writer, http.StatusOK, "You have been unsubscribed successfully.")
}

func (h *SubscribeHandler) Confirm(writer http.ResponseWriter, request *http.Request) {
	input, err := ParseAndValidateTokenInput(request)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	sub, err := repository.GetSubscriptionByToken(input.Token)
	if err != nil {
		response.RespondJSON(writer, http.StatusConflict, err.Error())
		return
	}

	if err := h.subService.Confirm(sub); err != nil {
		if err.Error() == "not found" {
			response.RespondJSON(writer, http.StatusNotFound, "Token not found")
		} else {
			response.RespondJSON(writer, http.StatusBadRequest, "Failed to confirm subscription: "+err.Error())
		}
		return
	}

	response.RespondJSON(writer, http.StatusOK, "Subscription confirmed successfully.")
}
