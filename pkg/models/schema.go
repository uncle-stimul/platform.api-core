package models

type LinkDefinition struct {
	ParentTable      string
	ParentField      string
	ParentEntity     string
	ChildrenTable    string
	ChildrenField    string
	ChildrenEntities []string
}

var SchemaSections = []Sections{
	{Module: "api-core", Endpoint: "/"},
	{Module: "api-core", Endpoint: "/settings/"},
	{Module: "api-core", Endpoint: "/settings/m/modules"},
	{Module: "api-core", Endpoint: "/settings/m/plugins"},
	{Module: "api-core", Endpoint: "/settings/a/users"},
	{Module: "api-core", Endpoint: "/settings/a/roles"},
	{Module: "api-core", Endpoint: "/settings/a/sections"},
	{Module: "api-core", Endpoint: "/settings/a/permissions"},
}

var SchemaPermissions = []Permissions{
	{Name: "mainpage-access", Description: "Доступ к главной странице платформы"},
	{Name: "settings-access", Description: "Доступ к основной странице с настройкам платформы"},
	{Name: "settings-monitoring-access", Description: "Доступ к страницам настроек с разделам мониторинга компонентов"},
	{Name: "settings-authentication-access", Description: "Доступ к страницам настроек с разделам управлением пользователями и доступом"},
}

var SchemaRoles = []Roles{
	{Name: "platform-admin", Description: "Администратор платформы"},
	{Name: "platform-user", Description: "Пользователь платформы"},
}

var SchemaUsers = []Users{{Username: "admin", Password: "admin", Status: true}}

var SchemaLinks = []LinkDefinition{
	{
		ParentTable: "permissions", ParentField: "Name", ParentEntity: "mainpage-access",
		ChildrenTable: "sections", ChildrenField: "Endpoint", ChildrenEntities: []string{"/"},
	},
	{
		ParentTable: "permissions", ParentField: "Name", ParentEntity: "settings-access",
		ChildrenTable: "sections", ChildrenField: "Endpoint", ChildrenEntities: []string{
			"/settings/m/modules", "/settings/m/plugins", "/settings/a/users",
			"/settings/a/roles", "/settings/a/sections", "/settings/a/permissions",
		},
	},
	{
		ParentTable: "roles", ParentField: "Name", ParentEntity: "platform-admin",
		ChildrenTable: "permissions", ChildrenField: "Name", ChildrenEntities: []string{
			"mainpage-access", "settings-access",
		},
	},
	{
		ParentTable: "roles", ParentField: "Name", ParentEntity: "platform-user",
		ChildrenTable: "permissions", ChildrenField: "Name", ChildrenEntities: []string{"mainpage-access"},
	},
	{
		ParentTable: "users", ParentField: "Username", ParentEntity: "admin",
		ChildrenTable: "roles", ChildrenField: "Name", ChildrenEntities: []string{"platform-admin"},
	},
}
