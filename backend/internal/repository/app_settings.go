package repository

import (
	"ithozyeva/database"
)

type AppSettingsRepository struct{}

func NewAppSettingsRepository() *AppSettingsRepository {
	return &AppSettingsRepository{}
}

// GetRaw возвращает JSONB-значение настройки в виде []byte. Возвращает
// (nil, nil), если ключа в таблице нет — вызывающий код применит default.
func (r *AppSettingsRepository) GetRaw(key string) ([]byte, error) {
	var raw []byte
	err := database.DB.Raw(
		`SELECT value::text FROM app_settings WHERE key = ?`,
		key,
	).Scan(&raw).Error
	return raw, err
}
