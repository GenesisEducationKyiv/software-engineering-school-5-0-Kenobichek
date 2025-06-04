package repository

import (
	"Weather-Forecast-API/internal/db"
	"Weather-Forecast-API/internal/models"
	"database/sql"
	"errors"
	"strings"
	"time"
)

func CreateSubscription(s *models.Subscription) error {
	s.NextNotifiedAt = time.Now().Add(time.Duration(s.FrequencyMinutes) * time.Minute)

	_, err := db.DB.Exec(`
		INSERT INTO subscriptions 
		(channel_type, channel_value, city, frequency_minutes, token, next_notified_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		s.ChannelType, s.ChannelValue, s.City,
		s.FrequencyMinutes, s.Token, s.NextNotifiedAt,
	)

	if err != nil && strings.Contains(err.Error(), "unique") {
		return errors.New("already subscribed")
	}
	return err
}

func ConfirmByToken(token string) error {
	result, err := db.DB.Exec(`
		UPDATE subscriptions
		SET confirmed = TRUE
		WHERE token = $1`, token)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("not found")
	}

	return nil
}

func UnsubscribeByToken(token string) error {
	result, err := db.DB.Exec(`DELETE FROM subscriptions WHERE token = $1`, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("not found")
	}
	return nil
}

func GetDueSubscriptions() []models.Subscription {
	rows, _ := db.DB.Query(`
		SELECT id, channel_type, channel_value, city, frequency_minutes
		FROM subscriptions
		WHERE confirmed = TRUE AND next_notified_at <= NOW()
	`)

	var subs []models.Subscription
	for rows.Next() {
		var s models.Subscription
		rows.Scan(&s.ID, &s.ChannelType, &s.ChannelValue, &s.City, &s.FrequencyMinutes)
		subs = append(subs, s)
	}

	return subs
}

func UpdateNextNotification(id int, next time.Time) {
	db.DB.Exec(`UPDATE subscriptions SET next_notified_at = $1 WHERE id = $2`, next, id)
}

func GetSubscriptionByToken(token string) (models.Subscription, error) {
	row := db.DB.QueryRow(`
		SELECT id, channel_type, channel_value, city, frequency_minutes, confirmed, token, next_notified_at, created_at
		FROM subscriptions
		WHERE token = $1
	`, token)

	var s models.Subscription
	err := row.Scan(
		&s.ID,
		&s.ChannelType,
		&s.ChannelValue,
		&s.City,
		&s.FrequencyMinutes,
		&s.Confirmed,
		&s.Token,
		&s.NextNotifiedAt,
		&s.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return s, errors.New("subscription not found")
	}
	if err != nil {
		return s, err
	}

	return s, nil
}
