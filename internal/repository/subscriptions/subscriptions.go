package subscriptions

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateSubscription(subscription *Info) error {
	_, err := r.db.Exec(`
		INSERT INTO subscriptions 
		(channel_type, channel_value, city, frequency_minutes, token, next_notified_at)
		VALUES ($1, $2, $3, $4, $5, NOW() + ($6 * interval '1 minute'))`,
		subscription.ChannelType, subscription.ChannelValue, subscription.City,
		subscription.FrequencyMinutes, subscription.Token, subscription.FrequencyMinutes,
	)

	if err != nil && strings.Contains(err.Error(), "unique") {
		return errors.New("already subscribed")
	}

	return nil
}
func (r *Repository) ConfirmByToken(token string) error {
	result, err := r.db.Exec(`
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

func (r *Repository) UnsubscribeByToken(token string) error {
	result, err := r.db.Exec(`DELETE FROM subscriptions WHERE token = $1`, token)
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

func (r *Repository) GetDueSubscriptions() []Info {
	rows, err := r.db.Query(
		`	SELECT id, channel_type, channel_value, city, frequency_minutes
			FROM subscriptions
			WHERE confirmed = TRUE AND next_notified_at <= NOW()`,
	)

	var subs []Info

	if err != nil {
		return subs
	}

	for rows.Next() {
		var s Info

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
func (r *Repository) UpdateNextNotification(id int, next time.Time) error {
	_, err := r.db.Exec(`UPDATE subscriptions SET next_notified_at = $1 WHERE id = $2`, next, id)
	if err != nil {
		return fmt.Errorf("failed to update next_notified_at for id %d: %w", id, err)
	}

	return nil
}

func (r *Repository) GetSubscriptionByToken(token string) (*Info, error) {
	row := r.db.QueryRow(`
		SELECT id, channel_type, channel_value, city, frequency_minutes, confirmed, token, next_notified_at, created_at
		FROM subscriptions
		WHERE token = $1
	`, token)

	var subscription Info
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
