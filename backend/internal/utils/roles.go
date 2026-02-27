package utils

import (
	"ithozyeva/internal/models"
)

func HasRole(roles []models.Role, role models.Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func RemoveRole(roles []models.Role, target models.Role) []models.Role {
	result := make([]models.Role, 0, len(roles))
	for _, r := range roles {
		if r != target {
			result = append(result, r)
		}
	}
	return result
}
