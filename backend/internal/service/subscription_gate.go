package service

import (
	"ithozyeva/config"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

// SubscriptionTierGate — единая точка решения «имеет ли пользователь доступ
// к функционалу, ограниченному tier'ом». Используется и middleware'ом
// RequireMinTier, и сервисными visibility-чекерами для комментов, чтобы
// поведение ровно совпадало.
//
// Поведение:
//   - SUBSCRIPTION_GATE_ENABLED=false  → всегда true (gate отключён глобально);
//   - роль ADMIN                       → всегда true (admin — универсальный
//     модератор, см. PR #316);
//   - active subscription tier         → level >= minLevel.
//
// Без этого helper'а middleware и checker могли разойтись и foreman
// получил бы доступ к комментам через /comments/:id/like при отключённом
// gate'е, но 403 на /ai-materials/:id/comments через middleware
// (или наоборот).
type SubscriptionTierGate struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionTierGate(repo *repository.SubscriptionRepository) *SubscriptionTierGate {
	return &SubscriptionTierGate{repo: repo}
}

// AllowsMinTier возвращает true, если member имеет право работать с
// функционалом, требующим tier'а уровня >= minLevel.
func (g *SubscriptionTierGate) AllowsMinTier(member *models.Member, minLevel int) bool {
	if !config.CFG.SubscriptionGateEnabled {
		return true
	}
	if member == nil {
		return false
	}
	for _, role := range member.Roles {
		if role == models.MemberRoleAdmin {
			return true
		}
	}
	level, ok := g.repo.GetUserEffectiveTierLevel(member.TelegramID)
	if !ok {
		return false
	}
	return level >= minLevel
}
