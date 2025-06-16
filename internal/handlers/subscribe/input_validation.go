package subscribe

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

type SubscribeInput struct {
	ChannelType  string
	ChannelValue string
	City         string
	Frequency    string
}

type TokenInput struct {
	Token string
}

func ParseAndValidateSubscribeInput(r *http.Request) (SubscribeInput, error) {
	channelValue := strings.TrimSpace(r.FormValue("email"))
	city := strings.TrimSpace(r.FormValue("city"))
	frequency := strings.TrimSpace(r.FormValue("frequency"))
	if channelValue == "" || city == "" || frequency == "" {
		return SubscribeInput{}, errors.New("invalid input")
	}
	channelType := strings.TrimSpace(r.FormValue("channelType"))
	if channelType == "" {
		channelType = "email"
	}
	if !IsValidChannel(channelType) {
		return SubscribeInput{}, errors.New("unsupported channelType")
	}
	return SubscribeInput{
		ChannelType:  channelType,
		ChannelValue: channelValue,
		City:         city,
		Frequency:    frequency,
	}, nil
}

func ParseAndValidateTokenInput(r *http.Request) (TokenInput, error) {
	token := strings.TrimSpace(chi.URLParam(r, "token"))

	if token == "" {
		return TokenInput{}, errors.New("invalid token")
	}
	return TokenInput{
		Token: token,
	}, nil
}

func IsValidChannel(channel string) bool {
	_, ok := SupportedChannels()[channel]

	return ok
}
