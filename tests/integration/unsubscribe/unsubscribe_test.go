package unsubscribe_test

import (
	"Weather-Forecast-API/internal/handlers/subscribe"
	"Weather-Forecast-API/tests/testdb"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"testing"
)

const (
	unsubscribeEndpoint = "/api/unsubscribe/"
	validToken          = "3fa85f64-5717-4562-b3fc-2c963f66afa6" //nolint:gosec
)

func TestUnsubscribeAPI(t *testing.T) {
	pg := testdb.New(t)
	database := pg.SQL

	_, err := database.Exec(
		`INSERT INTO subscriptions (channel_value, city, token, 
		frequency_minutes, channel_type, next_notified_at)
		VALUES ($1, $2, $3, $4, $5, NOW())`, "test@example.com", "Kyiv",
		validToken, 60, "email")
	require.NoError(t, err)

	testCases := []struct {
		name               string
		requestBody        map[string]string
		expectedStatusCode int
		expectUnsubscribe  bool
		expectedBody       string
	}{
		{
			name: "Successful unsubscribe",
			requestBody: map[string]string{
				"token": validToken,
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"` + subscribe.MsgUnsubscribedSuccess + `"}`,
			expectUnsubscribe:  true,
		},
		{
			name: "Invalid token",
			requestBody: map[string]string{
				"token": "invalid-token",
			},
			expectedStatusCode: http.StatusConflict,
			expectedBody:       `{"message":"` + subscribe.ErrTokenNotFound.Error() + `"}`,
			expectUnsubscribe:  false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dbMutex.Lock()
			defer dbMutex.Unlock()

			mockNotif := newMockNotificationService()
			appSrv := newTestServer(t, newAppRouterWithDB(t, database, mockNotif))

			_, contentType := multipartBody(t, tc.requestBody)

			reqURL := appSrv.URL + unsubscribeEndpoint + tc.requestBody["token"]
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
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
			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			log.Printf("response body: %s", string(bodyBytes))
			if tc.expectedBody != "" {
				require.True(t, json.Valid([]byte(tc.expectedBody)), "expectedBody is not valid JSON")
				require.True(t, json.Valid(bodyBytes), "response body is not valid JSON")
				assert.JSONEq(t, tc.expectedBody, string(bodyBytes), "response body mismatch")
			}
			if tc.expectUnsubscribe {
				require.Len(t, mockNotif.sentUnsubscribes, 1, "expected unsubscribe notification to be sent")
			} else {
				require.Len(t, mockNotif.sentUnsubscribes, 0, "did not expect unsubscribe notification to be sent")
			}
		})
	}
}
