package handler

import "ithozyeva/internal/models"

// sanitizePublicMember зануляет приватные поля Member перед отдачей в
// публичный API. TelegramID — PII: его утечка позволяла любому
// авторизованному пользователю эскалировать сессию до чужой через
// /auth/telegram/refresh.
func sanitizePublicMember(m *models.Member) {
	if m == nil {
		return
	}
	m.TelegramID = 0
}

// sanitizePublicMentor — то же для DTO ментора.
func sanitizePublicMentor(m *models.MentorModel) {
	if m == nil {
		return
	}
	m.TelegramID = 0
}
