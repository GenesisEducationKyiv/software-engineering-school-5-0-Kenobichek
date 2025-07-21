CREATE TABLE subscriptions (
	id SERIAL PRIMARY KEY,
	channel_type VARCHAR(50) NOT NULL CHECK (channel_type IN ('email', 'sms')),
	channel_value TEXT NOT NULL,
	city VARCHAR(100) NOT NULL,
	frequency_minutes INTEGER NOT NULL CHECK (frequency_minutes > 0),
	confirmed BOOLEAN NOT NULL DEFAULT FALSE,
	token UUID NOT NULL UNIQUE,
	next_notified_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	
	UNIQUE (channel_type, channel_value, city)
);