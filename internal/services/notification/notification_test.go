package notification

import (
	"errors"
	"testing"
)

type mockEmailNotifier struct {
	sendFunc func(to, message, subject string) error
}

func (m *mockEmailNotifier) Send(to, message, subject string) error {
	if m.sendFunc != nil {
		return m.sendFunc(to, message, subject)
	}
	return nil
}

func TestService_SendMessage(t *testing.T) {
	tests := []struct {
		name         string
		channelType  string
		channelValue string
		message      string
		subject      string
		notifierErr  error
		expectedErr  string
	}{
		{
			name:         "successful email send",
			channelType:  "email",
			channelValue: "test@example.com",
			message:      "Test message",
			subject:      "Test subject",
			notifierErr:  nil,
			expectedErr:  "",
		},
		{
			name:         "notifier error",
			channelType:  "email",
			channelValue: "test@example.com",
			message:      "Test message",
			subject:      "Test subject",
			notifierErr:  errors.New("failed to send"),
			expectedErr:  "failed to send",
		},
		{
			name:         "unsupported channel type",
			channelType:  "sms",
			channelValue: "+1234567890",
			message:      "Test message",
			subject:      "Test subject",
			notifierErr:  nil,
			expectedErr:  "unsupported channel type: sms",
		},
		{
			name:         "empty email address",
			channelType:  "email",
			channelValue: "",
			message:      "Test message",
			subject:      "Test subject",
			notifierErr:  nil,
			expectedErr:  "",
		},
		{
			name:         "empty message content",
			channelType:  "email",
			channelValue: "test@example.com",
			message:      "",
			subject:      "Test subject",
			notifierErr:  nil,
			expectedErr:  "",
		},
		{
			name:         "empty subject",
			channelType:  "email",
			channelValue: "test@example.com",
			message:      "Test message",
			subject:      "",
			notifierErr:  nil,
			expectedErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier := &mockEmailNotifier{
				sendFunc: func(to, message, subject string) error {
					return tt.notifierErr
				},
			}

			service := NewService(notifier)
			err := service.SendMessage(tt.channelType, tt.channelValue, tt.message, tt.subject)

			if tt.expectedErr == "" && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Errorf("expected error %s, got %v", tt.expectedErr, err)
				}
			}
		})
	}
}

func TestNewService(t *testing.T) {
	notifier := &mockEmailNotifier{}
	service := NewService(notifier)

	if service == nil {
		t.Error("expected service to be created, got nil")
	}
}

func TestService_SendMessage_NilNotifier(t *testing.T) {
	service := &Service{notifier: nil}

	err := service.SendMessage("email", "test@example.com", "message", "subject")

	if err == nil {
		t.Error("expected error when notifier is nil, got none")
	}
}