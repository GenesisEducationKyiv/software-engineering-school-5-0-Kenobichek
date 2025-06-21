package subscribe

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"Weather-Forecast-API/internal/repository"

	"github.com/go-chi/chi/v5"
)

// Mock implementations
type mockSubscriptionService struct {
	subscribeFunc   func(*repository.Subscription) error
	unsubscribeFunc func(*repository.Subscription) error
	confirmFunc     func(*repository.Subscription) error
}

func (m *mockSubscriptionService) Subscribe(sub *repository.Subscription) error {
	if m.subscribeFunc != nil {
		return m.subscribeFunc(sub)
	}
	return nil
}

func (m *mockSubscriptionService) Unsubscribe(sub *repository.Subscription) error {
	if m.unsubscribeFunc != nil {
		return m.unsubscribeFunc(sub)
	}
	return nil
}

func (m *mockSubscriptionService) Confirm(sub *repository.Subscription) error {
	if m.confirmFunc != nil {
		return m.confirmFunc(sub)
	}
	return nil
}

type mockNotificationService struct {
	sendMessageFunc func(channelType, channelValue, message, subject string) error
}

func (m *mockNotificationService) SendMessage(channelType, channelValue, message, subject string) error {
	if m.sendMessageFunc != nil {
		return m.sendMessageFunc(channelType, channelValue, message, subject)
	}
	return nil
}

func TestHandler_Subscribe(t *testing.T) {
	tests := []struct {
		name               string
		formData           url.Values
		subscriptionError  error
		notificationError  error
		expectedStatus     int
		expectedBody       string
	}{
		{
			name: "successful subscription",
			formData: url.Values{
				"email":     []string{"test@example.com"},
				"city":      []string{"London"},
				"frequency": []string{"daily"},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Subscription successful! Check your email to confirm."}`,
		},
		{
			name: "missing email",
			formData: url.Values{
				"city":      []string{"London"},
				"frequency": []string{"daily"},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"email is required"}`,
		},
		{
			name: "invalid frequency",
			formData: url.Values{
				"email":     []string{"test@example.com"},
				"city":      []string{"London"},
				"frequency": []string{"invalid"},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid frequency"}`,
		},
		{
			name: "subscription service error",
			formData: url.Values{
				"email":     []string{"test@example.com"},
				"city":      []string{"London"},
				"frequency": []string{"daily"},
			},
			subscriptionError: errors.New("already subscribed"),
			expectedStatus:    http.StatusConflict,
			expectedBody:      `{"message":"already subscribed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subService := &mockSubscriptionService{
				subscribeFunc: func(sub *repository.Subscription) error {
					return tt.subscriptionError
				},
			}
			notifService := &mockNotificationService{
				sendMessageFunc: func(channelType, channelValue, message, subject string) error {
					return tt.notificationError
				},
			}

			handler := NewHandler(subService, notifService)

			body := tt.formData.Encode()
			req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.Subscribe(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if strings.TrimSpace(w.Body.String()) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}

func TestHandler_Unsubscribe(t *testing.T) {
	tests := []struct {
		name              string
		token             string
		unsubscribeError  error
		expectedStatus    int
		expectedBody      string
	}{
		{
			name:           "successful unsubscribe",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Successfully unsubscribed"}`,
		},
		{
			name:           "missing token",
			token:          "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"token is required"}`,
		},
		{
			name:             "not found error",
			token:            "invalid-token",
			unsubscribeError: errors.New("not found"),
			expectedStatus:   http.StatusNotFound,
			expectedBody:     `{"message":"not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subService := &mockSubscriptionService{
				unsubscribeFunc: func(sub *repository.Subscription) error {
					return tt.unsubscribeError
				},
			}
			handler := NewHandler(subService, &mockNotificationService{})

			req := httptest.NewRequest(http.MethodGet, "/unsubscribe/"+tt.token, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("token", tt.token)
			req = req.WithContext(chi.RouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.Unsubscribe(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if strings.TrimSpace(w.Body.String()) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}

func TestHandler_Confirm(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		confirmError   error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful confirmation",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Subscription confirmed successfully"}`,
		},
		{
			name:           "missing token",
			token:          "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"token is required"}`,
		},
		{
			name:         "not found error",
			token:        "invalid-token",
			confirmError: errors.New("not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody: `{"message":"not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subService := &mockSubscriptionService{
				confirmFunc: func(sub *repository.Subscription) error {
					return tt.confirmError
				},
			}
			handler := NewHandler(subService, &mockNotificationService{})

			req := httptest.NewRequest(http.MethodGet, "/confirm/"+tt.token, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("token", tt.token)
			req = req.WithContext(chi.RouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.Confirm(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if strings.TrimSpace(w.Body.String()) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}