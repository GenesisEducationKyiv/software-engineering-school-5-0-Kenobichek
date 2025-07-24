package domain

import (
	"fmt"
	"sync"
)

type TemplateRepository struct {
	mu        sync.RWMutex
	templates map[string]*MessageTemplate
}

func (r *TemplateRepository) GetTemplateByName(name string) (*MessageTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tpl, ok := r.templates[name]
	if !ok {
		return nil, fmt.Errorf("template '%s' not found", name)
	}
	return tpl, nil
}

func NewTemplateRepository() *TemplateRepository {
	return &TemplateRepository{
		templates: map[string]*MessageTemplate{
			"confirm": {
				Subject: "Confirm your weather subscription",
				Message: "Hello! To confirm your subscription, please use the code: {{ .ConfirmToken }}",
			},
			"weather_update": {
				Subject: "Weather update for {{ .City }}",
				Message: "Current weather in {{ .City }}: {{ .Description }}. " +
					"Temperature: {{ .Temperature }}°C, Humidity: {{ .Humidity }}%.",
			},
			"unsubscribe": {
				Subject: "You have unsubscribed from weather alerts for {{ .City }}",
				Message: "You have successfully unsubscribed from weather notifications for {{ .City }}.",
			},
		},
	}
}
