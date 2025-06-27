package sendgridemailapi_test

import (
	"Weather-Forecast-API/external/sendgridemailapi"
	"errors"
	"net/http"
	"testing"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type mockSendgridClient struct {
	mockResponse *rest.Response
	mockError    error
}

func (m *mockSendgridClient) Send(email *mail.SGMailV3) (*rest.Response, error) {
	return m.mockResponse, m.mockError
}

func TestSendgridNotifier_Send(t *testing.T) {
	tests := []struct {
		name            string
		target          sendgridemailapi.NotificationTarget
		message         string
		subject         string
		mockResponse    *rest.Response
		mockError       error
		expectedError   bool
		expectedErrText string
	}{
		{
			name:          "Valid email sent successfully",
			target:        sendgridemailapi.NotificationTarget{Type: sendgridemailapi.ChannelEmail, Address: "test@example.com"},
			message:       "Test Message",
			subject:       "Test Subject",
			mockResponse:  &rest.Response{StatusCode: http.StatusOK},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:            "Invalid notification type",
			target:          sendgridemailapi.NotificationTarget{Type: "InvalidChannel", Address: "test@example.com"},
			message:         "Test Message",
			subject:         "Test Subject",
			mockResponse:    nil,
			mockError:       nil,
			expectedError:   true,
			expectedErrText: "invalid notification target type InvalidChannel, expected email",
		},
		{
			name: "Sendgrid error response",
			target: sendgridemailapi.NotificationTarget{
				Type:    sendgridemailapi.ChannelEmail,
				Address: "test@example.com",
			},
			message:         "Test Message",
			subject:         "Test Subject",
			mockResponse:    &rest.Response{StatusCode: http.StatusInternalServerError},
			mockError:       nil,
			expectedError:   true,
			expectedErrText: "sendgrid returned error status code 500",
		},
		{
			name: "Sendgrid client error",
			target: sendgridemailapi.NotificationTarget{
				Type:    sendgridemailapi.ChannelEmail,
				Address: "test@example.com",
			},
			message:         "Test Message",
			subject:         "Test Subject",
			mockResponse:    nil,
			mockError:       errors.New("client error"),
			expectedError:   true,
			expectedErrText: "failed to send email: client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockSendgridClient{
				mockResponse: tt.mockResponse,
				mockError:    tt.mockError,
			}
			notifier := sendgridemailapi.NewSendgridNotifier(mockClient, "TestSender", "sender@example.com")

			err := notifier.Send(tt.target, tt.message, tt.subject)
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
			if err != nil && tt.expectedErrText != "" && err.Error() != tt.expectedErrText {
				t.Errorf("expected error message: '%s', got: '%s'", tt.expectedErrText, err.Error())
			}
		})
	}
}
