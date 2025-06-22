package subscribe

import (
	"Weather-Forecast-API/internal/repository/subscriptions"
	"Weather-Forecast-API/internal/response"
	"strings"

	"github.com/google/uuid"
	"net/http"
)

type subscriptionManager interface {
	Subscribe(sub *subscriptions.Info) error
	Unsubscribe(sub *subscriptions.Info) error
	Confirm(sub *subscriptions.Info) error
	GetSubscriptionByToken(token string) (*subscriptions.Info, error)
}

type notificationManager interface {
	SendConfirmation(channel string, recipient string, token string) error
	SendUnsubscribe(channel string, recipient string, city string) error
}

type Handler struct {
	subscriptionService subscriptionManager
	notificationService notificationManager
}

func NewHandler(
	subscriptionService subscriptionManager,
	notificationService notificationManager,
) *Handler {
	return &Handler{
		subscriptionService: subscriptionService,
		notificationService: notificationService,
	}
}

func (h *Handler) Subscribe(writer http.ResponseWriter, request *http.Request) {
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

	token := uuid.NewString()

	sub := &subscriptions.Info{
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

	err = h.notificationService.SendConfirmation(input.ChannelType, input.ChannelValue, token)
	if err != nil {
		response.RespondJSON(writer, http.StatusInternalServerError, "Failed to send message. Error: "+err.Error())
		return
	}

	response.RespondJSON(writer, http.StatusOK, "Subscription successful. Confirmation sent.")
}

func (h *Handler) Unsubscribe(writer http.ResponseWriter, request *http.Request) {
	input, err := parseTokenFromRequest(request)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	sub, err := h.subscriptionService.GetSubscriptionByToken(input.Token)
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

	err = h.notificationService.SendUnsubscribe(sub.ChannelType, sub.ChannelValue, sub.City)
	if err != nil {
		response.RespondJSON(writer, http.StatusInternalServerError, "Failed to send message. Error: "+err.Error())
		return
	}

	response.RespondJSON(writer, http.StatusOK, "You have been unsubscribed successfully.")
}

func (h *Handler) Confirm(writer http.ResponseWriter, request *http.Request) {
	input, err := parseTokenFromRequest(request)
	if err != nil {
		response.RespondJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	sub, err := h.subscriptionService.GetSubscriptionByToken(input.Token)
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
