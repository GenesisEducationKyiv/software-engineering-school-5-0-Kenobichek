package repository

import (
	"errors"

	"Weather-Forecast-API/internal/db"
	"Weather-Forecast-API/internal/models"
)

func GetTemplateByName(name string) (*models.MessageTemplate, error) {
	var tpl models.MessageTemplate

	query := `SELECT subject, message FROM email_templates WHERE name = $1`

	err := db.DataBase.QueryRow(query, name).Scan(&tpl.Subject, &tpl.Message)
	if err != nil {
		return nil, errors.New("failed get template by name")
	}

	return &tpl, nil
}
