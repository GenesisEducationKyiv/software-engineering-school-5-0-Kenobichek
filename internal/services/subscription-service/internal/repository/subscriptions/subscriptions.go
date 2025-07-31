package subscriptions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type databaseManager interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Repository struct {
	db databaseManager
}

func New(db databaseManager) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSubscription(ctx context.Context, sub *Subscription) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO subscriptions 
		(channel_type, channel_value, city, frequency_minutes, token, next_notified_at)
		VALUES ($1, $2, $3, $4, $5, NOW() + ($6 * interval '1 minute'))`,
		sub.ChannelType, sub.ChannelValue, sub.City,
		sub.FrequencyMinutes, sub.Token, sub.FrequencyMinutes,
	)
	if err != nil && strings.Contains(err.Error(), "unique") {
		return errors.New("already subscribed")
	}
	return err
}

func (r *Repository) ConfirmByToken(ctx context.Context, token string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE subscriptions
		SET confirmed = TRUE
		WHERE token = $1`, token)
	if err != nil {
		return fmt.Errorf("failed to update subscription confirmation: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected after confirmation update: %w", err)
	}
	if rows == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

func (r *Repository) UnsubscribeByToken(ctx context.Context, token string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM subscriptions WHERE token = $1`, token)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected after deletion: %w", err)
	}
	if rows == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

func (r *Repository) GetDueSubscriptions(ctx context.Context) ([]Subscription, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, channel_type, channel_value, city, frequency_minutes
		 FROM subscriptions
		 WHERE confirmed = TRUE AND next_notified_at <= NOW()`,
	)
	var subs []Subscription
	if err != nil {
		return subs, fmt.Errorf("failed to get due subscriptions: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()
	
	for rows.Next() {
		var s Subscription
		if err := rows.Scan(&s.ID, &s.ChannelType, &s.ChannelValue, &s.City, &s.FrequencyMinutes); err != nil {
			return subs, fmt.Errorf("failed to scan due subscriptions: %w", err)
		}
		subs = append(subs, s)
	}
	if err = rows.Err(); err != nil {
		return subs, fmt.Errorf("failed to get due subscriptions: %w", err)
	}
	return subs, nil
}

func (r *Repository) UpdateNextNotification(ctx context.Context, id int64, next time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE subscriptions SET next_notified_at = $1 WHERE id = $2`, next, id)
	if err != nil {
		return fmt.Errorf("failed to update next_notified_at for id %d: %w", id, err)
	}
	return nil
}

func (r *Repository) GetSubscriptionByToken(ctx context.Context, token string) (*Subscription, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, channel_type, channel_value, city, frequency_minutes, confirmed, token, next_notified_at, created_at
		FROM subscriptions
		WHERE token = $1
	`, token)
	var sub Subscription
	err := row.Scan(
		&sub.ID,
		&sub.ChannelType,
		&sub.ChannelValue,
		&sub.City,
		&sub.FrequencyMinutes,
		&sub.Confirmed,
		&sub.Token,
		&sub.NextNotifiedAt,
		&sub.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("subscription not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription by token '%s': %w", token, err)
	}
	return &sub, nil
}
