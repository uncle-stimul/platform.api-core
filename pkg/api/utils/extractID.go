package utils

import "platform.api-core/pkg/models"

func ExtractRoleIDs(roles []models.Roles) []uint {
	ids := make([]uint, len(roles))
	for i, role := range roles {
		ids[i] = role.ID
	}
	return ids
}
