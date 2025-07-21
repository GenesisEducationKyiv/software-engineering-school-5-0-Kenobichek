package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"api-gateway/internal/kafka"

	"log"

	"github.com/go-chi/chi/v5"
)

type SubscribeHandler struct {
	Publisher *kafka.Publisher
}

func NewSubscribeHandler(publisher *kafka.Publisher) *SubscribeHandler {
	return &SubscribeHandler{Publisher: publisher}
}

func (h *SubscribeHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	city := r.FormValue("city")
	frequency := r.FormValue("frequency")

	if err := validateSubscriptionParams(email, city, frequency); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	frequencyMinutes, err := frequencyToMinutes(frequency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := SubscriptionCommand{
		Command:          "subscribe",
		ChannelType:      "email",
		ChannelValue:     email,
		City:             city,
		Frequency:        frequency,
		FrequencyMinutes: frequencyMinutes,
	}
	payload, err := json.Marshal(cmd)
	if err != nil {
		http.Error(w, "failed to marshal command", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := h.Publisher.Publish(ctx, email, payload); err != nil {
		log.Printf("failed to publish event: %v", err)
		http.Error(w, "failed to process subscription request", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write([]byte("Subscription event published")); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (h *SubscribeHandler) handleTokenCommand(
	w http.ResponseWriter, 
	r *http.Request, 
	command string, 
	validateFunc func(string) error,
	successMsg string,
) {
	token := chi.URLParam(r, "token")
	if err := validateFunc(token); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cmd := SubscriptionCommand{
		Command: command,
		Token:   token,
	}
	payload, err := json.Marshal(cmd)
	if err != nil {
		http.Error(w, "failed to marshal command", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := h.Publisher.Publish(ctx, token, payload); err != nil {
		http.Error(w, "failed to publish "+command+" event: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write([]byte(successMsg)); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (h *SubscribeHandler) ConfirmSubscription(w http.ResponseWriter, r *http.Request) {
	h.handleTokenCommand(w, r, "confirm", validateConfirmSubscriptionParams, "Confirm event published")
}

func (h *SubscribeHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	h.handleTokenCommand(w, r, "unsubscribe", validateUnsubscribeParams, "Unsubscribe event published")
}
