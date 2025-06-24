package subscribe_test

import (
	"Weather-Forecast-API/internal/handlers/subscribe"
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/repository/subscriptions"
	"Weather-Forecast-API/internal/routes"
	"Weather-Forecast-API/internal/services/subscription"
	"Weather-Forecast-API/tests/testdb"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"mime/multipart"

	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type notificationManager interface {
	SendConfirmation(channel string, recipient string, token string) error
	SendUnsubscribe(channel string, recipient string, city string) error
}

var dbMutex sync.Mutex

type mockNotificationService struct {
	sentConfirmations []confirmationCall
	sentUnsubscribes  []unsubscribeCall
	mutex             sync.Mutex
}

type confirmationCall struct {
	channel   string
	recipient string
	token     string
}

type unsubscribeCall struct {
	channel   string
	recipient string
	city      string
}

func newMockNotificationService() *mockNotificationService {
	return &mockNotificationService{
		sentConfirmations: make([]confirmationCall, 0),
		sentUnsubscribes:  make([]unsubscribeCall, 0),
	}
}

func (m *mockNotificationService) SendConfirmation(channel, recipient, token string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sentConfirmations = append(m.sentConfirmations, confirmationCall{
		channel:   channel,
		recipient: recipient,
		token:     token,
	})
	return nil
}

func (m *mockNotificationService) SendUnsubscribe(channel, recipient, city string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sentUnsubscribes = append(m.sentUnsubscribes, unsubscribeCall{
		channel:   channel,
		recipient: recipient,
		city:      city,
	})
	return nil
}

func (m *mockNotificationService) SendMessage() error { return nil }

func (m *mockNotificationService) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sentConfirmations = m.sentConfirmations[:0]
	m.sentUnsubscribes = m.sentUnsubscribes[:0]
}

func newAppRouterWithDB(t *testing.T, db *sql.DB, notifSvc notificationManager) http.Handler {
	t.Helper()

	subsRepo := subscriptions.New(db)
	subsSvc := subscription.NewService(subsRepo)
	subscribeHandler := subscribe.NewHandler(subsSvc, notifSvc)

	weatherHandler := weather.NewHandler(nil, 5*time.Second)

	router := routes.NewHTTPRouter()
	routes.NewService(weatherHandler, subscribeHandler, router).RegisterRoutes()

	return router
}

func TestSubscribeAPI(t *testing.T) {
	pg := testdb.New(t)
	database := pg.SQL

	testCases := []struct {
		name               string
		requestBody        map[string]string
		expectedStatusCode int
		expectConfirmation bool
		expectedBody       string
	}{
		{
			name: "Successful subscription",
			requestBody: map[string]string{
				"email":     "test@example.com",
				"city":      "Kyiv",
				"frequency": "daily",
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"` + subscribe.MsgSubscriptionSuccess + `"}`,
			expectConfirmation: true,
		},
		{
			name: "Missing email",
			requestBody: map[string]string{
				"city":      "Kyiv",
				"frequency": "daily",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"` + subscribe.ErrInvalidInput.Error() + `"}`,
			expectConfirmation: false,
		},
		{
			name: "Already subscribed",
			requestBody: map[string]string{
				"email":     "test@example.com",
				"city":      "Kyiv",
				"frequency": "daily",
			},
			expectedBody:       `{"message":"` + subscribe.ErrAlreadySubscribed.Error() + `"}`,
			expectedStatusCode: http.StatusConflict,
			expectConfirmation: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dbMutex.Lock()
			defer dbMutex.Unlock()

			mockNotif := newMockNotificationService()
			appSrv := newTestServer(t, newAppRouterWithDB(t, database, mockNotif))

			body, contentType := multipartBody(t, tc.requestBody)

			reqURL := appSrv.URL + "/api/subscribe"
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, body)
			require.NoError(t, err)

			req.Header.Set("Content-Type", contentType)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Println("failed to close response body")
					return
				}
			}()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			bodyBytes, _ := io.ReadAll(resp.Body)

			log.Printf("response body: %s", string(bodyBytes))
			if tc.expectedBody != "" {
				require.True(t, json.Valid([]byte(tc.expectedBody)), "expectedBody is not valid JSON")
				require.True(t, json.Valid(bodyBytes), "response body is not valid JSON")
				assert.JSONEq(t, tc.expectedBody, string(bodyBytes), "response body mismatch")
			}
			if tc.expectConfirmation {
				require.Len(t, mockNotif.sentConfirmations, 1,
					"expected confirmation to be sent")
			} else {
				require.Len(t, mockNotif.sentConfirmations, 0,
					"did not expect confirmation to be sent")
			}
		})
	}
}

func newTestServer(t *testing.T, h http.Handler) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv
}

func multipartBody(t *testing.T, fields map[string]string) (*bytes.Buffer, string) {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for k, v := range fields {
		require.NoError(t, writer.WriteField(k, v))
	}

	require.NoError(t, writer.Close())
	return &buf, writer.FormDataContentType()
}
