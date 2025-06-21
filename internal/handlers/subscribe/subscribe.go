package subscribe

import (
	"Weather-Forecast-API/internal/repository"
	"Weather-Forecast-API/internal/response"
	"Weather-Forecast-API/internal/services/notification"
	"Weather-Forecast-API/internal/services/subscription"
	"strings"

	"github.com/google/uuid"
	"net/http"
)

type SubscriptionManager interface {
	Subscribe(writer http.ResponseWriter, request *http.Request)
	Unsubscribe(writer http.ResponseWriter, request *http.Request)
	Confirm(writer http.ResponseWriter, request *http.Request)
}

type SubscriptionHandler struct {
	subscriptionService subscription.SubscriptionService
	notificationService notification.NotificationService
}

func NewSubscribeHandler(
	subscriptionService subscription.SubscriptionService,
	notificationService notification.NotificationService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		notificationService: notificationService,
	}
}

func (h *SubscriptionHandler) Subscribe(writer http.ResponseWriter, request *http.Request) {
	input, err := parseSubscribeInput(request)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	frequencyMinutes, err := convertFrequencyToMinutes(input.Frequency)
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

	sub := &repository.Subscription{
		ChannelType:      input.ChannelType,
		ChannelValue:     input.ChannelValue,
		City:             input.City,
		FrequencyMinutes: frequencyMinutes,
		Token:            token,
	}

	if err := h.subscriptionService.Subscribe(sub); err != nil {
		response.RespondJSON(writer, http.StatusConflict, "Already subscribed or DB error")
		return
	}

	message := strings.ReplaceAll(template.Message, "{{ confirm_token }}", token)

	err = h.notificationService.SendMessage(input.ChannelType, input.ChannelValue, message, template.Subject)
	if err != nil {
		response.RespondJSON(writer, http.StatusInternalServerError, "Failed to send message. Error: "+err.Error())
		return
	}

	response.RespondJSON(writer, http.StatusOK, "Subscription successful. Confirmation sent.")
}

func (h *SubscriptionHandler) Unsubscribe(writer http.ResponseWriter, request *http.Request) {
	input, err := parseTokenFromRequest(request)
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

	if err := h.subscriptionService.Unsubscribe(sub); err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.RespondJSON(writer, http.StatusNotFound, "Token not found")
		} else {
			response.RespondJSON(writer, http.StatusBadRequest, "Failed to confirm subscription: "+err.Error())
		}
		return
	}

	message := strings.ReplaceAll(template.Message, "{{ city }}", sub.City)

	err = h.notificationService.SendMessage(sub.ChannelType, sub.ChannelValue, message, template.Subject)
	if err != nil {
		response.RespondJSON(writer, http.StatusInternalServerError, "Failed to send message. Error: "+err.Error())
		return
	}

	response.RespondJSON(writer, http.StatusOK, "You have been unsubscribed successfully.")
}

func (h *SubscriptionHandler) Confirm(writer http.ResponseWriter, request *http.Request) {
	input, err := parseTokenFromRequest(request)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	sub, err := repository.GetSubscriptionByToken(input.Token)
	if err != nil {
		response.RespondJSON(writer, http.StatusConflict, err.Error())
		return
	}

	if err := h.subscriptionService.Confirm(sub); err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.RespondJSON(writer, http.StatusNotFound, "Token not found")
		} else {
			response.RespondJSON(writer, http.StatusBadRequest, "Failed to confirm subscription: "+err.Error())
		}
		return
	}

	response.RespondJSON(writer, http.StatusOK, "Subscription confirmed successfully.")
}
