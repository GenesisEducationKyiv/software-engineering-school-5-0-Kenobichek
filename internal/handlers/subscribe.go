package handlers

import (
	"Weather-Forecast-API/internal/models"
	"Weather-Forecast-API/internal/repository"
	notificationService "Weather-Forecast-API/internal/services/notification"
	subscriptionService "Weather-Forecast-API/internal/services/subscription"
	"strings"

	"Weather-Forecast-API/internal/utilities"
	"github.com/google/uuid"
	"net/http"
)

type SubscribeHandler struct {
	subscriptionService subscriptionService.SubscriptionService
	notificationService notificationService.NotificationService
}

func NewSubscribeHandler(
	subscriptionService subscriptionService.SubscriptionService,
	notificationService notificationService.NotificationService) *SubscribeHandler {
	return &SubscribeHandler{
		subscriptionService: subscriptionService,
		notificationService: notificationService,
	}
}

func (h *SubscribeHandler) Subscribe(writer http.ResponseWriter, request *http.Request) {
	input, err := utilities.ParseAndValidateSubscribeInput(request)
	if err != nil {
		utilities.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	frequencyMinutes, err := utilities.ConvertFrequency(input.Frequency)
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
		ChannelType:      input.ChannelType,
		ChannelValue:     input.ChannelValue,
		City:             input.City,
		FrequencyMinutes: frequencyMinutes,
		Token:            token,
	}

	if err := h.subscriptionService.Subscribe(sub); err != nil {
		utilities.RespondJSON(writer, http.StatusConflict, "Already subscribed or DB error")
		return
	}

	message := strings.ReplaceAll(template.Message, "{{ confirm_token }}", token)

	err = h.notificationService.SendMessage(input.ChannelType, input.ChannelValue, message, template.Subject)
	if err != nil {
		utilities.RespondJSON(writer, http.StatusInternalServerError, "Failed to send message. Error: "+err.Error())
		return
	}

	utilities.RespondJSON(writer, http.StatusOK, "Subscription successful. Confirmation sent.")
}

func (h *SubscribeHandler) Unsubscribe(writer http.ResponseWriter, request *http.Request) {
	input, err := utilities.ParseAndValidateTokenInput(request)
	if err != nil {
		utilities.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	template, err := repository.GetTemplateByName("unsubscribe")
	if err != nil {
		utilities.RespondJSON(writer, http.StatusInternalServerError, "Failed to load unsubscribe template")
		return
	}

	sub, err := repository.GetSubscriptionByToken(input.Token)
	if err != nil {
		utilities.RespondJSON(writer, http.StatusConflict, err.Error())
		return
	}

	if err := h.subscriptionService.Unsubscribe(sub); err != nil {
		if err.Error() == "not found" {
			utilities.RespondJSON(writer, http.StatusNotFound, "Token not found")
		} else {
			utilities.RespondJSON(writer, http.StatusBadRequest, "Failed to confirm subscription: "+err.Error())
		}
		return
	}

	message := strings.ReplaceAll(template.Message, "{{ city }}", sub.City)

	err = h.notificationService.SendMessage(sub.ChannelType, sub.ChannelValue, message, template.Subject)
	if err != nil {
		utilities.RespondJSON(writer, http.StatusInternalServerError, "Failed to send message. Error: "+err.Error())
		return
	}

	utilities.RespondJSON(writer, http.StatusOK, "You have been unsubscribed successfully.")
}

func (h *SubscribeHandler) Confirm(writer http.ResponseWriter, request *http.Request) {
	input, err := utilities.ParseAndValidateTokenInput(request)
	if err != nil {
		utilities.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	sub, err := repository.GetSubscriptionByToken(input.Token)
	if err != nil {
		utilities.RespondJSON(writer, http.StatusConflict, err.Error())
		return
	}

	if err := h.subscriptionService.Confirm(sub); err != nil {
		if err.Error() == "not found" {
			utilities.RespondJSON(writer, http.StatusNotFound, "Token not found")
		} else {
			utilities.RespondJSON(writer, http.StatusBadRequest, "Failed to confirm subscription: "+err.Error())
		}
		return
	}

	utilities.RespondJSON(writer, http.StatusOK, "Subscription confirmed successfully.")
}
