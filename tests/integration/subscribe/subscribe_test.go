package subscribe_test

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
	subscribeEndpoint = "/api/subscribe"
)

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

			reqURL := appSrv.URL + subscribeEndpoint
			log.Println("Request URL:", reqURL)
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
			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

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
