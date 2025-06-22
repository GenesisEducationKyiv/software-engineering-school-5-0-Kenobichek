package subscriptions_test

import (
	"Weather-Forecast-API/internal/repository/subscriptions"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
	"time"
)

//nolint:gocyclo
func TestRepository(t *testing.T) { //nolint:gocognit
	db, mock, _ := sqlmock.New()
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("failed to close database connection: %v", closeErr)
		}
	}()

	repo := subscriptions.New(db)

	t.Run("CreateSubscription", func(t *testing.T) {
		tests := []struct {
			name        string
			input       *subscriptions.Info
			mockSetup   func()
			expectedErr error
		}{
			{
				name: "success",
				input: &subscriptions.Info{
					ChannelType:      "email",
					ChannelValue:     "test@example.com",
					City:             "New York",
					FrequencyMinutes: 60,
					Token:            "test-token",
				},
				mockSetup: func() {
					mock.ExpectExec("^INSERT INTO subscriptions").
						WithArgs("email", "test@example.com", "New York", 60, "test-token", 60).
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
				expectedErr: nil,
			},
			{
				name: "already subscribed",
				input: &subscriptions.Info{
					ChannelType:      "email",
					ChannelValue:     "duplicate@example.com",
					City:             "London",
					FrequencyMinutes: 30,
					Token:            "duplicate-token",
				},
				mockSetup: func() {
					mock.ExpectExec("^INSERT INTO subscriptions").
						WillReturnError(errors.New("unique violation"))
				},
				expectedErr: errors.New("already subscribed"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.mockSetup()
				err := repo.CreateSubscription(tt.input)
				if (err != nil && tt.expectedErr == nil) ||
					(err == nil && tt.expectedErr != nil) ||
					(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			})
		}
	})

	t.Run("ByToken", func(t *testing.T) {
		tests := []struct {
			name        string
			token       string
			mockSetup   func()
			expectedErr error
		}{
			{
				name:  "success update",
				token: "valid-token",
				mockSetup: func() {
					mock.ExpectExec("^UPDATE subscriptions SET confirmed = TRUE WHERE token = \\$1").
						WithArgs("valid-token").
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
				expectedErr: nil,
			},
			{
				name:  "not found",
				token: "missing-token",
				mockSetup: func() {
					mock.ExpectExec("^UPDATE subscriptions SET confirmed = TRUE WHERE token = \\$1").
						WithArgs("missing-token").
						WillReturnResult(sqlmock.NewResult(1, 0))
				},
				expectedErr: errors.New("not found"),
			},
			{
				name:  "success delete",
				token: "valid-token",
				mockSetup: func() {
					mock.ExpectExec("^DELETE FROM subscriptions WHERE token = \\$1").
						WithArgs("valid-token").
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
				expectedErr: nil,
			},
			{
				name:  "not found",
				token: "missing-token",
				mockSetup: func() {
					mock.ExpectExec("^DELETE FROM subscriptions WHERE token = \\$1").
						WithArgs("missing-token").
						WillReturnResult(sqlmock.NewResult(1, 0))
				},
				expectedErr: errors.New("not found"),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.mockSetup()
				err := repo.ConfirmByToken(tt.token)
				if (err != nil && tt.expectedErr == nil) ||
					(err == nil && tt.expectedErr != nil) ||
					(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			})
		}
	})

	t.Run("GetDueSubscriptions", func(t *testing.T) {
		tests := []struct {
			name         string
			mockSetup    func()
			expectedSubs []subscriptions.Info
		}{
			{
				name: "has due subscriptions",
				mockSetup: func() {
					rows := sqlmock.NewRows([]string{"id", "channel_type", "channel_value", "city", "frequency_minutes"}).
						AddRow(1, "email", "test1@example.com", "New York", 60).
						AddRow(2, "sms", "1234567890", "London", 30)
					mock.ExpectQuery("^SELECT id, channel_type, channel_value, city, frequency_minutes FROM subscriptions").
						WillReturnRows(rows)
				},
				expectedSubs: []subscriptions.Info{
					{ID: 1, ChannelType: "email", ChannelValue: "test1@example.com", City: "New York", FrequencyMinutes: 60},
					{ID: 2, ChannelType: "sms", ChannelValue: "1234567890", City: "London", FrequencyMinutes: 30},
				},
			},
			{
				name: "no due subscriptions",
				mockSetup: func() {
					mock.ExpectQuery("^SELECT id, channel_type, channel_value, city, frequency_minutes FROM subscriptions").
						WillReturnRows(sqlmock.NewRows([]string{"id", "channel_type", "channel_value", "city", "frequency_minutes"}))
				},
				expectedSubs: []subscriptions.Info{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.mockSetup()
				subs := repo.GetDueSubscriptions()
				if len(subs) != len(tt.expectedSubs) {
					t.Errorf("expected %d subscriptions, got %d", len(tt.expectedSubs), len(subs))
				}
			})
		}
	})

	t.Run("UpdateNextNotification", func(t *testing.T) {
		tests := []struct {
			name        string
			id          int
			next        time.Time
			mockSetup   func()
			expectedErr error
		}{
			{
				name: "success",
				id:   1,
				next: time.Now().Add(60 * time.Minute),
				mockSetup: func() {
					mock.ExpectExec("^UPDATE subscriptions SET next_notified_at").
						WithArgs(sqlmock.AnyArg(), 1).
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
				expectedErr: nil,
			},
			{
				name: "failure",
				id:   2,
				next: time.Now().Add(15 * time.Minute),
				mockSetup: func() {
					mock.ExpectExec("^UPDATE subscriptions SET next_notified_at").
						WithArgs(sqlmock.AnyArg(), 2).
						WillReturnError(errors.New("update failure"))
				},
				expectedErr: errors.New("failed to update next_notified_at for id 2: update failure"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.mockSetup()
				err := repo.UpdateNextNotification(tt.id, tt.next)
				if (err != nil && tt.expectedErr == nil) ||
					(err == nil && tt.expectedErr != nil) ||
					(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			})
		}
	})

	t.Run("GetSubscriptionByToken", func(t *testing.T) {
		tests := []struct {
			name        string
			token       string
			mockSetup   func()
			expectedSub *subscriptions.Info
			expectedErr error
		}{
			{
				name:  "found subscription",
				token: "valid-token",
				mockSetup: func() {
					rows := sqlmock.NewRows([]string{
						"id", "channel_type", "channel_value", "city", "frequency_minutes",
						"confirmed", "token", "next_notified_at", "created_at"}).
						AddRow(1, "email", "test@example.com", "New York", 60,
							true, "valid-token", time.Now(), time.Now())
					mock.ExpectQuery("^SELECT id, channel_type, channel_value").
						WithArgs("valid-token").
						WillReturnRows(rows)
				},
				expectedSub: &subscriptions.Info{
					ID:               1,
					ChannelType:      "email",
					ChannelValue:     "test@example.com",
					City:             "New York",
					FrequencyMinutes: 60,
					Confirmed:        true,
					Token:            "valid-token",
					NextNotifiedAt:   time.Now(),
					CreatedAt:        time.Now(),
				},
				expectedErr: nil,
			},
			{
				name:  "not found",
				token: "missing-token",
				mockSetup: func() {
					mock.ExpectQuery("^SELECT id, channel_type, channel_value").
						WithArgs("missing-token").
						WillReturnError(sql.ErrNoRows)
				},
				expectedSub: &subscriptions.Info{},
				expectedErr: errors.New("subscription not found"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.mockSetup()
				sub, err := repo.GetSubscriptionByToken(tt.token)
				if (err != nil && tt.expectedErr == nil) ||
					(err == nil && tt.expectedErr != nil) ||
					(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				if tt.expectedSub != nil && sub.Token != tt.expectedSub.Token {
					t.Errorf("expected subscription %v, got %v", tt.expectedSub, sub)
				}
			})
		}
	})
}
