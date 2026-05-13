package handler

import "ithozyeva/internal/models"

// sanitizePublicMember зануляет приватные поля Member перед отдачей в
// публичный API. PII-поля:
//   - TelegramID: утечка позволяла любому авторизованному пользователю
//     эскалировать сессию до чужой через /auth/telegram/refresh.
//   - Birthday: дата рождения — sensitive personal data (152-ФЗ / GDPR),
//     не должна светиться в публичных списках /api/members или в карточке
//     участника, открытой любому залогиненному. У владельца профиля своя
//     дата возвращается через /me (sanitize здесь не применяется).
func sanitizePublicMember(m *models.Member) {
	if m == nil {
		return
	}
	m.TelegramID = 0
	m.Birthday = nil
}

// sanitizePublicMentor — то же для DTO ментора.
func sanitizePublicMentor(m *models.MentorModel) {
	if m == nil {
		return
	}
	m.TelegramID = 0
	m.Birthday = nil
}
