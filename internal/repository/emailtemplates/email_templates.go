package emailtemplates

import (
	"Weather-Forecast-API/internal/templates"
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetTemplateByName(name templates.Name) (*MessageTemplate, error) {
	var tpl MessageTemplate

	query := `SELECT subject, message FROM email_templates WHERE name = $1`

	err := r.db.QueryRow(query, name).Scan(&tpl.Subject, &tpl.Message)
	if err != nil {
		return nil, errors.New("failed get template by name")
	}

	return &tpl, nil
}
