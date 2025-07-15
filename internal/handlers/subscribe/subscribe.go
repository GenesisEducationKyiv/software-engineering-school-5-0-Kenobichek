package subscribe

import (
	"Weather-Forecast-API/internal/response"
	"errors"

	"net/http"
)

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrAlreadySubscribed   = errors.New("already subscribed or db error")
	ErrTokenNotFound       = errors.New("token not found")
	ErrFailedToSendMessage = errors.New("failed to send message")
	ErrFailedToConfirm     = errors.New("failed to confirm subscription")

	MsgSubscriptionSuccess   = "subscription successful. confirmation sent."
	MsgUnsubscribedSuccess   = "you have been unsubscribed successfully."
	MsgSubscriptionConfirmed = "subscription confirmed successfully."
)

type subscriptionManager interface {
}

type notificationManager interface {
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
	// input, err := parseSubscribeInput(request)
	// if err != nil {
	// 	response.RespondJSON(writer, http.StatusBadRequest, ErrInvalidInput.Error())
	// 	return
	// }

	// frequencyMinutes, err := convertFrequencyToMinutes(input.Frequency)
	// if err != nil {
	// 	response.RespondJSON(writer, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// token := uuid.NewString()

	// sub := &subscriptions.Info{
	// 	ChannelType:      input.ChannelType,
	// 	ChannelValue:     input.ChannelValue,
	// 	City:             input.City,
	// 	FrequencyMinutes: frequencyMinutes,
	// 	Token:            token,
	// }

	// if err := h.subscriptionService.Subscribe(sub); err != nil {
	// 	if errors.Is(err, ErrAlreadySubscribed) || strings.Contains(err.Error(), "already subscribed") {
	// 		response.RespondJSON(writer, http.StatusConflict, ErrAlreadySubscribed.Error())
	// 	} else {
	// 		response.RespondJSON(writer, http.StatusConflict, err.Error())
	// 	}
	// 	return
	// }

	// err = h.notificationService.SendConfirmation(input.ChannelType, input.ChannelValue, token)
	// if err != nil {
	// 	response.RespondJSON(writer, http.StatusInternalServerError, ErrFailedToSendMessage.Error()+". Error: "+err.Error())
	// 	return
	// }

	response.RespondJSON(writer, http.StatusOK, MsgSubscriptionSuccess)
}

func (h *Handler) Unsubscribe(writer http.ResponseWriter, request *http.Request) {
	// input, err := parseTokenFromRequest(request)
	// if err != nil {
	// 	response.RespondJSON(writer, http.StatusBadRequest, ErrInvalidInput.Error())
	// 	return
	// }

	// sub, err := h.subscriptionService.GetSubscriptionByToken(input.Token)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "invalid input syntax for type uuid") ||
	// 		errors.Is(err, ErrTokenNotFound) ||
	// 		strings.Contains(err.Error(), "not found") {
	// 		response.RespondJSON(writer, http.StatusConflict, ErrTokenNotFound.Error())
	// 	} else {
	// 		response.RespondJSON(writer, http.StatusConflict, err.Error())
	// 	}
	// 	return
	// }

	// if err := h.subscriptionService.Unsubscribe(sub); err != nil {
	// 	if errors.Is(err, ErrTokenNotFound) || strings.Contains(err.Error(), "not found") {
	// 		response.RespondJSON(writer, http.StatusNotFound, ErrTokenNotFound.Error())
	// 	} else {
	// 		response.RespondJSON(writer, http.StatusBadRequest, ErrFailedToConfirm.Error()+": "+err.Error())
	// 	}
	// 	return
	// }

	// err = h.notificationService.SendUnsubscribe(sub.ChannelType, sub.ChannelValue, sub.City)
	// if err != nil {
	// 	response.RespondJSON(writer, http.StatusInternalServerError, ErrFailedToSendMessage.Error()+". Error: "+err.Error())
	// 	return
	// }

	response.RespondJSON(writer, http.StatusOK, MsgUnsubscribedSuccess)
}

func (h *Handler) Confirm(writer http.ResponseWriter, request *http.Request) {
	// input, err := parseTokenFromRequest(request)
	// if err != nil {
	// 	response.RespondJSON(writer, http.StatusBadRequest, ErrInvalidInput.Error())
	// 	return
	// }

	// sub, err := h.subscriptionService.GetSubscriptionByToken(input.Token)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "invalid input syntax for type uuid") ||
	// 		errors.Is(err, ErrTokenNotFound) ||
	// 		strings.Contains(err.Error(), "not found") {
	// 		response.RespondJSON(writer, http.StatusConflict, ErrTokenNotFound.Error())
	// 	} else {
	// 		response.RespondJSON(writer, http.StatusConflict, err.Error())
	// 	}
	// 	return
	// }

	// if err := h.subscriptionService.Confirm(sub); err != nil {
	// 	if errors.Is(err, ErrTokenNotFound) || strings.Contains(err.Error(), "not found") {
	// 		response.RespondJSON(writer, http.StatusNotFound, ErrTokenNotFound.Error())
	// 	} else {
	// 		response.RespondJSON(writer, http.StatusBadRequest, ErrFailedToConfirm.Error()+": "+err.Error())
	// 	}
	// 	return
	// }

	response.RespondJSON(writer, http.StatusOK, MsgSubscriptionConfirmed)
}
