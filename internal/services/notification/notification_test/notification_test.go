package notification_test

import (
	"Weather-Forecast-API/internal/services/notification"
	"errors"
	"testing"

	"Weather-Forecast-API/external/sendgridemailapi"
)

type mockEmailNotifier struct {
	sendFunc func(to, message, subject string) error
}

func (m *mockEmailNotifier) Send(to, message, subject string) error {
	return m.sendFunc(to, message, subject)
}

func TestService_SendMessage(t *testing.T) {
	tests := []struct {
		name          string
		channelType   string
		channelValue  string
		message       string
		subject       string
		sendFunc      func(to, message, subject string) error
		expectedError string
	}{
		{
			name:         "valid email channel",
			channelType:  string(sendgridemailapi.ChannelEmail),
			channelValue: "test@example.com",
			message:      "Hello World",
			subject:      "Test Subject",
			sendFunc: func(to, message, subject string) error {
				if to == "test@example.com" && message == "Hello World" && subject == "Test Subject" {
					return nil
				}
				return errors.New("unexpected parameters")
			},
			expectedError: "",
		},
		{
			name:         "unsupported channel type",
			channelType:  "sms",
			channelValue: "1234567890",
			message:      "Hello World",
			subject:      "Test Subject",
			sendFunc: func(to, message, subject string) error {
				return nil
			},
			expectedError: "unsupported channel type: sms",
		},
		{
			name:         "send returns error",
			channelType:  string(sendgridemailapi.ChannelEmail),
			channelValue: "test@example.com",
			message:      "Hello World",
			subject:      "Test Subject",
			sendFunc: func(to, message, subject string) error {
				return errors.New("failed to send email")
			},
			expectedError: "failed to send email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier := &mockEmailNotifier{sendFunc: tt.sendFunc}
			service := notification.NewService(notifier)

			err := service.SendMessage(tt.channelType, tt.channelValue, tt.message, tt.subject)

			if tt.expectedError == "" && err != nil {
				t.Errorf("expected no error but got: %v", err)
			} else if tt.expectedError != "" && (err == nil || err.Error() != tt.expectedError) {
				t.Errorf("expected error %q but got: %v", tt.expectedError, err)
			}
		})
	}
}
