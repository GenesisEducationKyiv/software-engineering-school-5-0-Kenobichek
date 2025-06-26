package emailtemplates_test

import (
	"Weather-Forecast-API/internal/repository/emailtemplates"
	"Weather-Forecast-API/internal/templates"
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_GetTemplateByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("failed to close database connection: %v", closeErr)
		}
	}()

	repo := emailtemplates.New(db)

	tests := []struct {
		name          string
		templateName  templates.Name
		mockQuery     func()
		expectedTpl   *emailtemplates.MessageTemplate
		expectedError error
	}{
		{
			name:         "valid template",
			templateName: templates.Name("welcome"),
			mockQuery: func() {
				mock.ExpectQuery(`SELECT subject, message FROM email_templates WHERE name = \$1`).
					WithArgs("welcome").
					WillReturnRows(sqlmock.NewRows([]string{"subject", "message"}).
						AddRow("Welcome!", "Hello, welcome to our service."))
			},
			expectedTpl: &emailtemplates.MessageTemplate{
				Subject: "Welcome!",
				Message: "Hello, welcome to our service.",
			},
			expectedError: nil,
		},
		{
			name:         "template not found",
			templateName: templates.Name("goodbye"),
			mockQuery: func() {
				mock.ExpectQuery(`SELECT subject, message FROM email_templates WHERE name = \$1`).
					WithArgs("goodbye").
					WillReturnError(sql.ErrNoRows)
			},
			expectedTpl:   nil,
			expectedError: errors.New("failed get template by name"),
		},
		{
			name:         "database error",
			templateName: templates.Name("error_case"),
			mockQuery: func() {
				mock.ExpectQuery(`SELECT subject, message FROM email_templates WHERE name = \$1`).
					WithArgs("error_case").
					WillReturnError(errors.New("random database error"))
			},
			expectedTpl:   nil,
			expectedError: errors.New("failed get template by name"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()

			tpl, err := repo.GetTemplateByName(tt.templateName)

			if (err != nil && tt.expectedError == nil) ||
				(err == nil && tt.expectedError != nil) ||
				(err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tpl != nil && tt.expectedTpl != nil {
				if tpl.Subject != tt.expectedTpl.Subject || tpl.Message != tt.expectedTpl.Message {
					t.Errorf("expected template %+v, got %+v", tt.expectedTpl, tpl)
				}
			}

			if tpl == nil && tt.expectedTpl != nil {
				t.Errorf("expected template %+v, got nil", tt.expectedTpl)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled mock expectations: %s", err)
	}
}
