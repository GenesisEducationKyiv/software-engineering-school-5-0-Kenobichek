package domain

import "fmt"

type TemplateRepository struct {
	templates map[string]*MessageTemplate
}

func NewTemplateRepository() *TemplateRepository {
	return &TemplateRepository{
		templates: map[string]*MessageTemplate{
			"confirm": {
				Subject: "Confirm your weather subscription",
				Message: "Hello! To confirm your subscription, please use the code: {{ confirm_token }}",
			},
			"weather_update": {
				Subject: "Weather update for {{ city }}",
				Message: "Current weather in {{ city }}: {{ description }}. Temperature: {{ temperature }}°C, Humidity: {{ humidity }}%.",
			},
			"unsubscribe": {
				Subject: "You have unsubscribed from weather alerts for {{ city }}",
				Message: "You have successfully unsubscribed from weather notifications for {{ city }}.",
			},
		},
	}
}

func (r *TemplateRepository) GetTemplateByName(name string) (*MessageTemplate, error) {
	tpl, ok := r.templates[name]
	if !ok {
		return nil, fmt.Errorf("template '%s' not found", name)
	}
	return tpl, nil
}
