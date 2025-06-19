package repository

import (
	"errors"

	"Weather-Forecast-API/internal/db"
)

func GetTemplateByName(name string) (*MessageTemplate, error) {
	var tpl MessageTemplate

	query := `SELECT subject, message FROM email_templates WHERE name = $1`

	err := db.DataBase.QueryRow(query, name).Scan(&tpl.Subject, &tpl.Message)
	if err != nil {
		return nil, errors.New("failed get template by name")
	}

	return &tpl, nil
}
