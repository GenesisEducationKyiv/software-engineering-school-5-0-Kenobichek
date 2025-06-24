package weather_test

import (
	"Weather-Forecast-API/external/openweather"
	"Weather-Forecast-API/internal/handlers/subscribe"
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/repository/subscriptions"
	"Weather-Forecast-API/internal/routes"
	"Weather-Forecast-API/internal/weatherprovider/openweatherprovider"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	mockAPIKey = "TEST_API_KEY"
	testCity   = "Kyiv"
)

type stubSubSvc struct{}

func (stubSubSvc) Subscribe(*subscriptions.Info) error                        { return nil }
func (stubSubSvc) Unsubscribe(*subscriptions.Info) error                      { return nil }
func (stubSubSvc) Confirm(*subscriptions.Info) error                          { return nil }
func (stubSubSvc) GetSubscriptionByToken(string) (*subscriptions.Info, error) { return nil, nil }
func (stubNotifSvc) SendConfirmation(string, string, string) error {
	return nil
}
func (stubNotifSvc) SendUnsubscribe(string, string, string) error { return nil }

type stubNotifSvc struct{}

func (stubNotifSvc) SendMessage() error { return nil }

type openWeatherAPIMock struct {
	geoSrv     *httptest.Server
	weatherSrv *httptest.Server
}

func newOpenWeatherAPIMock(
	t *testing.T,
	city string,
	geoResponse any, geoStatus int,
	weatherResponse any, weatherStatus int,
) *openWeatherAPIMock {
	t.Helper()

	geoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if city != "" {
			assert.Contains(t, r.URL.RawQuery, "q="+city)
		}
		assert.Contains(t, r.URL.RawQuery, "appid="+mockAPIKey)
		w.WriteHeader(geoStatus)
		if s, ok := geoResponse.(string); ok {
			_, _ = io.WriteString(w, s)
			return
		}
		_ = json.NewEncoder(w).Encode(geoResponse)
	}))
	weatherSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.RawQuery, "appid="+mockAPIKey)
		w.WriteHeader(weatherStatus)
		if s, ok := weatherResponse.(string); ok {
			_, _ = io.WriteString(w, s)
			return
		}
		_ = json.NewEncoder(w).Encode(weatherResponse)
	}))

	mock := &openWeatherAPIMock{geoSrv: geoSrv, weatherSrv: weatherSrv}
	t.Cleanup(mock.Close)

	return mock
}

func (m *openWeatherAPIMock) Close() {
	m.geoSrv.Close()
	m.weatherSrv.Close()
}

func newAppRouter(t *testing.T, owMock *openWeatherAPIMock) http.Handler {
	t.Helper()

	httpClient := &http.Client{Timeout: 2 * time.Second}

	geoSvc := openweather.NewGeocodingService(httpClient, owMock.geoSrv.URL, mockAPIKey)
	owAPI := openweather.NewOpenWeatherAPI(httpClient, owMock.weatherSrv.URL, mockAPIKey)
	weatherProv := openweatherprovider.NewOpenWeatherProvider(geoSvc, owAPI)

	weatherHandler := weather.NewHandler(weatherProv, 5*time.Second)
	subscribeHandler := subscribe.NewHandler(stubSubSvc{}, stubNotifSvc{})

	router := routes.NewHTTPRouter()
	routes.NewService(weatherHandler, subscribeHandler, router).RegisterRoutes()

	return router
}

func newTestServer(t *testing.T, h http.Handler) *httptest.Server {
	t.Helper()

	srv := httptest.NewUnstartedServer(h)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv.Listener = ln
	srv.Start()

	t.Cleanup(srv.Close)
	return srv
}

// ─────────────────────────────────────────────────────────────────────
// Tests
// ─────────────────────────────────────────────────────────────────────

func TestWeatherAPI(t *testing.T) {
	geoRespSuccess := []map[string]any{{"lat": 50.45, "lon": 30.52}}
	weatherRespSuccess := map[string]any{
		"main":    map[string]any{"temp": 23.1, "humidity": 60.0},
		"weather": []map[string]any{{"description": "clear sky"}},
	}
	errorResp := map[string]any{"message": "boom"}

	testCases := []struct {
		name               string
		city               string
		geoResponse        any
		geoStatus          int
		weatherResponse    any
		weatherStatus      int
		expectedStatusCode int
	}{
		{
			name:               "Success",
			city:               testCity,
			geoResponse:        geoRespSuccess,
			geoStatus:          http.StatusOK,
			weatherResponse:    weatherRespSuccess,
			weatherStatus:      http.StatusOK,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "City Not Found",
			city:               "NoCity",
			geoResponse:        []any{},
			geoStatus:          http.StatusOK,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Geocoding API Error",
			city:               testCity,
			geoResponse:        errorResp,
			geoStatus:          http.StatusInternalServerError,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Weather API Error",
			city:               testCity,
			geoResponse:        geoRespSuccess,
			geoStatus:          http.StatusOK,
			weatherResponse:    errorResp,
			weatherStatus:      http.StatusInternalServerError,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Geocoding API Invalid JSON",
			city:               testCity,
			geoResponse:        `{"lat: 50.45}`,
			geoStatus:          http.StatusOK,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Weather API Invalid JSON",
			city:               testCity,
			geoResponse:        geoRespSuccess,
			geoStatus:          http.StatusOK,
			weatherResponse:    `{"main": {"temp": 23.1}}`,
			weatherStatus:      http.StatusOK,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			owMock := newOpenWeatherAPIMock(t, tc.city, tc.geoResponse, tc.geoStatus, tc.weatherResponse, tc.weatherStatus)
			appSrv := newTestServer(t, newAppRouter(t, owMock))

			reqURL := appSrv.URL + "/api/weather?city=" + tc.city
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
			require.NoError(t, err)
			resp, err := appSrv.Client().Do(req)
			require.NoError(t, err)
			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Println("failed to close response body")
					return
				}
			}()
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedStatusCode == http.StatusOK {
				var res struct {
					Temperature float64 `json:"temperature"`
					Humidity    float64 `json:"humidity"`
					Description string  `json:"description"`
				}
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&res))
				assert.InDelta(t, 23.1, res.Temperature, 0.01)
				assert.Equal(t, 60.0, res.Humidity)
				assert.Equal(t, "clear sky", res.Description)
			} else {
				var errBody struct {
					Message string `json:"message"`
				}
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &errBody))
				assert.NotEmpty(t, errBody.Message)
			}
		})
	}
}

func TestWeather_Timeout(t *testing.T) {
	geoResp := []map[string]any{{"lat": 50.45, "lon": 30.52}}
	sleepSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(3 * time.Second)
	}))
	defer sleepSrv.Close()

	owMock := &openWeatherAPIMock{
		geoSrv: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(geoResp)
		})),
		weatherSrv: sleepSrv,
	}
	t.Cleanup(owMock.Close)

	appSrv := newTestServer(t, newAppRouter(t, owMock))

	reqURL := appSrv.URL + "/api/weather?city=" + testCity
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	require.NoError(t, err)
	resp, err := appSrv.Client().Do(req)
	require.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
			return
		}
	}()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
