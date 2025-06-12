INSERT INTO email_templates (name, subject, message)
VALUES
('confirm', 'Confirm your subscription', 'Use this token to confirm your subscription: {{ confirm_token }}'),
('unsubscribe', 'You have unsubscribed', 'You have been unsubscribed from weather updates for {{ city }}'),
('weather_update', 'Weather Update for {{ city }}', 'Weather in {{ city }}: {{ description }}, {{ temperature }}°C, humidity: {{ humidity }}%')
