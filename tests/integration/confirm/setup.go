package confirm_test

import (
	"Weather-Forecast-API/internal/handlers/subscribe"
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/repository/subscriptions"
	"Weather-Forecast-API/internal/routes"
	"Weather-Forecast-API/internal/services/subscription"
	"bytes"
	"database/sql"
	"github.com/stretchr/testify/require"
	"mime/multipart"
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
