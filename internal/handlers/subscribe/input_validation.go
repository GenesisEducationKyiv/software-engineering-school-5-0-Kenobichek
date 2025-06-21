package subscribe

import (
	"Weather-Forecast-API/external/sendgridemailapi"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

const (
	hourlyMinutes = 60
	dailyMinutes  = 1440
)

func parseSubscribeInput(r *http.Request) (Input, error) {
	channelValue := strings.TrimSpace(r.FormValue("email"))
	city := strings.TrimSpace(r.FormValue("city"))
	frequency := strings.TrimSpace(r.FormValue("frequency"))
	if channelValue == "" || city == "" || frequency == "" {
		return Input{}, errors.New("invalid input")
	}
	channelType := strings.TrimSpace(r.FormValue("channelType"))
	if channelType == "" {
		channelType = string(sendgridemailapi.ChannelEmail)
	}
	if !isValidChannel(channelType) {
		return Input{}, errors.New("unsupported channelType")
	}
	return Input{
		ChannelType:  channelType,
		ChannelValue: channelValue,
		City:         city,
		Frequency:    frequency,
	}, nil
}

func parseTokenFromRequest(r *http.Request) (TokenInput, error) {
	token := strings.TrimSpace(chi.URLParam(r, "token"))

	if token == "" {
		return TokenInput{}, errors.New("invalid token")
	}
	return TokenInput{
		Token: token,
	}, nil
}

func isValidChannel(channel string) bool {
	_, ok := validChannels()[channel]

	return ok
}

func validChannels() map[string]struct{} {
	return map[string]struct{}{
		"email": {},
		// "sms":   {},
	}
}

func frequencyToMinutes() map[string]int {
	return map[string]int{
		"hourly": hourlyMinutes,
		"daily":  dailyMinutes,
	}
}
