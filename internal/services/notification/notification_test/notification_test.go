package notification_test

import (
	"Weather-Forecast-API/internal/services/notification"
	"errors"
	"testing"

	"Weather-Forecast-API/external/sendgridemailapi"
	"Weather-Forecast-API/internal/handlers/weather"
	"Weather-Forecast-API/internal/repository/emailtemplates"
	"Weather-Forecast-API/internal/templates"
)

type mockNotifier struct {
	sendErr error
}

func (m *mockNotifier) Send(to, message, subject string) error {
	return m.sendErr
}

type mockTemplateRepository struct {
	templates map[templates.Name]*emailtemplates.MessageTemplate
	getErr    error
}

func (m *mockTemplateRepository) GetTemplateByName(name templates.Name) (*emailtemplates.MessageTemplate, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	t, exists := m.templates[name]
	if !exists {
		return nil, errors.New("template not found")
	}
	return t, nil
}

func TestService_SendConfirmation(t *testing.T) {
	tests := []struct {
		name        string
		channel     string
		recipient   string
		token       string
		templates   map[templates.Name]*emailtemplates.MessageTemplate
		templateErr error
		notifierErr error
		wantErr     bool
	}{
		{
			name:      "valid email confirmation",
			channel:   string(sendgridemailapi.ChannelEmail),
			recipient: "test@example.com",
			token:     "12345",
			templates: map[templates.Name]*emailtemplates.MessageTemplate{
				templates.Confirm: {
					Message: "Confirmation token: {{ confirm_token }}",
					Subject: "Confirm your action",
				},
			},
			wantErr: false,
		},
		{
			name:        "unsupported channel",
			channel:     "sms",
			recipient:   "test@example.com",
			token:       "12345",
			templates:   nil,
			templateErr: nil,
			wantErr:     true,
		},
		{
			name:        "error getting template",
			channel:     string(sendgridemailapi.ChannelEmail),
			recipient:   "test@example.com",
			token:       "12345",
			templateErr: errors.New("template error"),
			wantErr:     true,
		},
		{
			name:      "error sending email",
			channel:   string(sendgridemailapi.ChannelEmail),
			recipient: "test@example.com",
			token:     "12345",
			templates: map[templates.Name]*emailtemplates.MessageTemplate{
				templates.Confirm: {
					Message: "Confirmation token: {{ confirm_token }}",
					Subject: "Confirm your action",
				},
			},
			notifierErr: errors.New("send error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier := &mockNotifier{sendErr: tt.notifierErr}
			templateRepo := &mockTemplateRepository{
				templates: tt.templates,
				getErr:    tt.templateErr,
			}
			service := notification.NewService(notifier, templateRepo)

			err := service.SendConfirmation(tt.channel, tt.recipient, tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendConfirmation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_SendWeatherUpdate(t *testing.T) {
	tests := []struct {
		name        string
		channel     string
		recipient   string
		metrics     weather.Metrics
		templates   map[templates.Name]*emailtemplates.MessageTemplate
		templateErr error
		notifierErr error
		wantErr     bool
	}{
		{
			name:      "valid weather update",
			channel:   string(sendgridemailapi.ChannelEmail),
			recipient: "test@example.com",
			metrics: weather.Metrics{
				City:        "TestCity",
				Description: "Sunny",
				Temperature: 25.5,
				Humidity:    60,
			},
			templates: map[templates.Name]*emailtemplates.MessageTemplate{
				templates.WeatherUpdate: {
					Message: "Weather in {{ city }}: {{ description }}, {{ temperature }}°C, {{ humidity }}% humidity.",
					Subject: "Weather Update for {{ city }}",
				},
			},
			wantErr: false,
		},
		{
			name:      "error getting weather template",
			channel:   string(sendgridemailapi.ChannelEmail),
			recipient: "test@example.com",
			metrics: weather.Metrics{
				City:        "TestCity",
				Description: "Sunny",
				Temperature: 25.5,
				Humidity:    60,
			},
			templateErr: errors.New("template error"),
			wantErr:     true,
		},
		{
			name:        "error sending weather update",
			channel:     string(sendgridemailapi.ChannelEmail),
			recipient:   "test@example.com",
			metrics:     weather.Metrics{},
			notifierErr: errors.New("send error"),
			wantErr:     true,
		},
		{
			name:      "unsupported weather update channel",
			channel:   "sms",
			recipient: "test@example.com",
			metrics:   weather.Metrics{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier := &mockNotifier{sendErr: tt.notifierErr}
			templateRepo := &mockTemplateRepository{
				templates: tt.templates,
				getErr:    tt.templateErr,
			}
			service := notification.NewService(notifier, templateRepo)

			err := service.SendWeatherUpdate(tt.channel, tt.recipient, tt.metrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendWeatherUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_SendUnsubscribe(t *testing.T) {
	tests := []struct {
		name        string
		channel     string
		recipient   string
		city        string
		templates   map[templates.Name]*emailtemplates.MessageTemplate
		templateErr error
		notifierErr error
		wantErr     bool
	}{
		{
			name:      "valid unsubscribe email",
			channel:   string(sendgridemailapi.ChannelEmail),
			recipient: "test@example.com",
			city:      "TestCity",
			templates: map[templates.Name]*emailtemplates.MessageTemplate{
				templates.Unsubscribe: {
					Message: "You have unsubscribed from updates for {{ city }}.",
					Subject: "Unsubscribe Confirmation",
				},
			},
			wantErr: false,
		},
		{
			name:        "error getting unsubscribe template",
			channel:     string(sendgridemailapi.ChannelEmail),
			recipient:   "test@example.com",
			city:        "TestCity",
			templateErr: errors.New("template error"),
			wantErr:     true,
		},
		{
			name:        "error sending unsubscribe email",
			channel:     string(sendgridemailapi.ChannelEmail),
			recipient:   "test@example.com",
			city:        "TestCity",
			notifierErr: errors.New("send error"),
			wantErr:     true,
		},
		{
			name:      "unsupported unsubscribe channel",
			channel:   "sms",
			recipient: "test@example.com",
			city:      "TestCity",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier := &mockNotifier{sendErr: tt.notifierErr}
			templateRepo := &mockTemplateRepository{
				templates: tt.templates,
				getErr:    tt.templateErr,
			}
			service := notification.NewService(notifier, templateRepo)

			err := service.SendUnsubscribe(tt.channel, tt.recipient, tt.city)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendUnsubscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
