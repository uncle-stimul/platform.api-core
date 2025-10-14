package utils

import "platform.api-core/pkg/models"

func ExtractRolesNames(roles []models.Roles) []string {
	names := make([]string, len(roles))
	for i, role := range roles {
		names[i] = role.Name
	}
	return names
}

func ExtractPermissionsNames(permissions []models.Permissions) []string {
	names := make([]string, len(permissions))
	for i, permission := range permissions {
		names[i] = permission.Name
	}
	return names
}

func ExtractSectionsEndpoints(sections []models.Sections) []string {
	endpoints := make([]string, len(sections))
	for i, section := range sections {
		endpoints[i] = section.Endpoint
	}
	return endpoints
}
