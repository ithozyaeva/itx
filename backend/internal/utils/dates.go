package utils

import (
	"errors"
	"ithozyeva/internal/models"
	"time"
)

func ParseDate(dateStr *string) (*models.DateOnly, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}

	parsedDate, err := time.Parse(models.DateFormat, *dateStr)
	if err != nil {
		return nil, errors.New("invalid date format (expected YYYY-MM-DD)")
	}
	return models.NewDateOnly(&parsedDate), nil
}
