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

type SubscriptionCommand struct {
	Command          string `json:"command"`
	ChannelType      string `json:"channel_type,omitempty"`
	ChannelValue     string `json:"channel_value,omitempty"`
	City             string `json:"city,omitempty"`
	Frequency        string `json:"frequency,omitempty"`
	FrequencyMinutes int    `json:"frequency_minutes,omitempty"`
	Token            string `json:"token,omitempty"`
}

type SubscribeHandler struct {
	Publisher *kafka.Publisher
}

func NewSubscribeHandler(publisher *kafka.Publisher) *SubscribeHandler {
	return &SubscribeHandler{Publisher: publisher}
}

func (h *SubscribeHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	city := r.FormValue("city")
	frequency := r.FormValue("frequency")

	var frequencyMinutes int
	switch frequency {
	case "hourly":
		frequencyMinutes = 60
	case "daily":
		frequencyMinutes = 1440
	default:
		http.Error(w, "invalid frequency", http.StatusBadRequest)
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
	payload, _ := json.Marshal(cmd)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := h.Publisher.Publish(ctx, email, payload); err != nil {
		http.Error(w, "failed to publish event: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write([]byte("Subscription event published")); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (h *SubscribeHandler) ConfirmSubscription(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}
	cmd := SubscriptionCommand{
		Command: "confirm",
		Token:   token,
	}
	payload, _ := json.Marshal(cmd)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := h.Publisher.Publish(ctx, token, payload); err != nil {
		http.Error(w, "failed to publish confirm event: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write([]byte("Confirm event published")); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (h *SubscribeHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}
	cmd := SubscriptionCommand{
		Command: "unsubscribe",
		Token:   token,
	}
	payload, _ := json.Marshal(cmd)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := h.Publisher.Publish(ctx, token, payload); err != nil {
		http.Error(w, "failed to publish unsubscribe event: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write([]byte("Unsubscribe event published")); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
