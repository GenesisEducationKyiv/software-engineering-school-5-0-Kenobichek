package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"Weather-Forecast-API/internal/db"
)

func CreateSubscription(subscription *Subscription) error {
	subscription.NextNotifiedAt = time.Now().Add(time.Duration(subscription.FrequencyMinutes) * time.Minute)

	_, err := db.DataBase.Exec(`
		INSERT INTO subscriptions 
		(channel_type, channel_value, city, frequency_minutes, token, next_notified_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		subscription.ChannelType, subscription.ChannelValue, subscription.City,
		subscription.FrequencyMinutes, subscription.Token, subscription.NextNotifiedAt,
	)

	if err != nil && strings.Contains(err.Error(), "unique") {
		return errors.New("already subscribed")
	}

	return nil
}

func ConfirmByToken(token string) error {
	result, err := db.DataBase.Exec(`
		UPDATE subscriptions
		SET confirmed = TRUE
		WHERE token = $1`, token)
	if err != nil {
		return errors.New("failed update subscription confirmation")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed rows affected")
	}

	if rows == 0 {
		return errors.New("not found")
	}

	return nil
}

func UnsubscribeByToken(token string) error {
	result, err := db.DataBase.Exec(`DELETE FROM subscriptions WHERE token = $1`, token)
	if err != nil {
		return errors.New("failed delete subscription")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed rows affected")
	}

	if rows == 0 {
		return errors.New("not found")
	}

	return nil
}

func GetDueSubscriptions() []Subscription {
	rows, err := db.DataBase.Query(
		`	SELECT id, channel_type, channel_value, city, frequency_minutes
			FROM subscriptions
			WHERE confirmed = TRUE AND next_notified_at <= NOW()`,
	)

	var subs []Subscription

	if err != nil {
		return subs
	}

	for rows.Next() {
		var s Subscription

		if err := rows.Scan(&s.ID, &s.ChannelType, &s.ChannelValue, &s.City, &s.FrequencyMinutes); err != nil {
			return subs
		}

		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return subs
	}

	return subs
}
func UpdateNextNotification(id int, next time.Time) error {
	_, err := db.DataBase.Exec(`UPDATE subscriptions SET next_notified_at = $1 WHERE id = $2`, next, id)
	if err != nil {
		return fmt.Errorf("failed to update next_notified_at for id %d: %w", id, err)
	}

	return nil
}

func GetSubscriptionByToken(token string) (*Subscription, error) {
	row := db.DataBase.QueryRow(`
		SELECT id, channel_type, channel_value, city, frequency_minutes, confirmed, token, next_notified_at, created_at
		FROM subscriptions
		WHERE token = $1
	`, token)

	var subscription Subscription
	err := row.Scan(
		&subscription.ID,
		&subscription.ChannelType,
		&subscription.ChannelValue,
		&subscription.City,
		&subscription.FrequencyMinutes,
		&subscription.Confirmed,
		&subscription.Token,
		&subscription.NextNotifiedAt,
		&subscription.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return &subscription, errors.New("subscription not found")
	}

	if err != nil {
		return &subscription, fmt.Errorf("failed to get subscription by token '%s': %w", token, err)
	}

	return &subscription, nil
}
