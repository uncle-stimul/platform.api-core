package db

import (
	"platform.api-core/pkg/db/methods"
	"platform.api-core/pkg/db/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func initShema(pgdb *gorm.DB, log *logrus.Logger) error {
	log.Debug("Инициализация схемы модуля api-platform")
	if err := pgdb.AutoMigrate(
		&models.Users{},
		&models.Roles{},
		&models.Permissions{},
		&models.Sections{},
	); err != nil {
		return err
	}

	mainSection := models.Sections{Module: "api-platform", URL: "/main"}
	stat, err := methods.Create[models.Sections](pgdb, map[string]interface{}{"url": "/main"}, &mainSection)
	if err != nil {
		log.WithError(err).Fatalf("При создании раздела \"%s\" возникла не предвиденная ошибка:", mainSection.URL)
	}
	if stat {
		log.Debugf("Объект \"%s\" успешно создан в таблице \"sections\"", mainSection.URL)
	} else {
		log.Debugf("Объект \"%s\" уже существует в таблице \"sections\"", mainSection.URL)
	}

	settingSection := models.Sections{Module: "api-platform", URL: "/settings/status"}
	stat, err = methods.Create[models.Sections](pgdb, map[string]interface{}{"url": "/settings/status"}, &settingSection)
	if err != nil {
		log.WithError(err).Fatalf("При создании раздела \"%s\" возникла не предвиденная ошибка:", settingSection.URL)
	}
	if stat {
		log.Debugf("Объект \"%s\" успешно создан в таблице \"sections\"", settingSection.URL)
	} else {
		log.Debugf("Объект \"%s\" уже существует в таблице \"sections\"", settingSection.URL)
	}

	userSettingSection := models.Sections{Module: "api-platform", URL: "/settings/users"}
	stat, err = methods.Create[models.Sections](pgdb, map[string]interface{}{"url": "/settings/users"}, &userSettingSection)
	if err != nil {
		log.WithError(err).Fatalf("При создании раздела \"%s\" возникла не предвиденная ошибка:", userSettingSection.URL)
	}
	if stat {
		log.Debugf("Объект \"%s\" успешно создан в таблице \"sections\"", userSettingSection.URL)
	} else {
		log.Debugf("Объект \"%s\" уже существует в таблице \"sections\"", userSettingSection.URL)
	}

	roleSettingSection := models.Sections{Module: "api-platform", URL: "/settings/roles"}
	stat, err = methods.Create[models.Sections](pgdb, map[string]interface{}{"url": "/settings/roles"}, &roleSettingSection)
	if err != nil {
		log.WithError(err).Fatalf("При создании раздела \"%s\" возникла не предвиденная ошибка:", roleSettingSection.URL)
	}
	if stat {
		log.Debugf("Объект \"%s\" успешно создан в таблице \"sections\"", roleSettingSection.URL)
	} else {
		log.Debugf("Объект \"%s\" уже существует в таблице \"sections\"", roleSettingSection.URL)
	}

	permissionSettingSection := models.Sections{Module: "api-platform", URL: "/settings/permissions"}
	stat, err = methods.Create[models.Sections](pgdb, map[string]interface{}{"url": "/settings/permissions"}, &permissionSettingSection)
	if err != nil {
		log.WithError(err).Fatalf("При создании раздела \"%s\" возникла не предвиденная ошибка:", permissionSettingSection.URL)
	}
	if stat {
		log.Debugf("Объект \"%s\" успешно создан в таблице \"sections\"", permissionSettingSection.URL)
	} else {
		log.Debugf("Объект \"%s\" уже существует в таблице \"sections\"", permissionSettingSection.URL)
	}

	mainpagePermission := models.Permissions{Name: "mainpage-access", Description: "Доступ к главной странице платформы"}
	stat, err = methods.Create[models.Permissions](pgdb, map[string]interface{}{"name": "mainpage-access"}, &mainpagePermission)
	if err != nil {
		log.WithError(err).Fatalf("При создании объекта \"%s\" возникла не предвиденная ошибка:", mainpagePermission.Name)
	}
	if stat {
		log.Debugf("Объект \"%s\" успешно создан в таблице \"permissions\"", mainpagePermission.Name)
		mainpagePermissionSections := []models.Sections{{ID: 1}}
		err = methods.SetAssociations[models.Permissions, models.Sections](pgdb, &mainpagePermission, "Sections", mainpagePermissionSections)
		if err != nil {
			log.WithError(err).Fatal("При создании встроенных связей возникла не предвиденная ошибка:")
		}
	} else {
		log.Debugf("Объект \"%s\" уже существует в таблице \"permissions\"", mainpagePermission.Name)
	}

	settingsPermission := models.Permissions{Name: "settings-access", Description: "Доступ к настройкам платформы"}
	stat, err = methods.Create[models.Permissions](pgdb, map[string]interface{}{"name": "settings-access"}, &settingsPermission)
	if err != nil {
		log.WithError(err).Fatalf("При создании объекта \"%s\" возникла не предвиденная ошибка:", settingsPermission.Name)
	}
	if stat {
		log.Debugf("Объект \"%s\" успешно создан в таблице \"permissions\"", settingsPermission.Name)
		settingsPermissionSections := []models.Sections{{ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}
		err = methods.SetAssociations[models.Permissions, models.Sections](pgdb, &settingsPermission, "Sections", settingsPermissionSections)
		if err != nil {
			log.WithError(err).Fatal("При создании встроенных связей возникла не предвиденная ошибка:")
		}
	} else {
		log.Debugf("Объект \"%s\" уже существует в таблице \"permissions\"", settingsPermission.Name)
	}

	adminRole := models.Roles{Name: "platform-admin", Description: "Администратор платформы"}
	stat, err = methods.Create[models.Roles](pgdb, map[string]interface{}{"name": "platform-admin"}, &adminRole)
	if err != nil {
		log.WithError(err).Fatalf("При создании роли \"%s\" возникла не предвиденная ошибка:", adminRole.Name)
	}
	if stat {
		log.Debugf("Объект \"%s\" успешно создан в таблице \"roles\"", adminRole.Name)
		adminRolePermissions := []models.Permissions{{ID: 1}, {ID: 2}}
		err = methods.SetAssociations[models.Roles, models.Permissions](pgdb, &adminRole, "Permissions", adminRolePermissions)
		if err != nil {
			log.WithError(err).Fatal("При создании встроенных связей возникла не предвиденная ошибка:")
		}
	} else {
		log.Debugf("Объект \"%s\" уже существует в таблице \"roles\"", adminRole.Name)
	}

	userRole := models.Roles{Name: "platform-user", Description: "Пользователь платформы"}
	stat, err = methods.Create[models.Roles](pgdb, map[string]interface{}{"name": "platform-user"}, &userRole)
	if err != nil {
		log.WithError(err).Fatalf("При создании роли \"%s\" возникла не предвиденная ошибка:", userRole.Name)
	}
	if stat {
		log.Debugf("Объект \"%s\" успешно создан в таблице \"roles\"", userRole.Name)
		userRolePermissions := []models.Permissions{{ID: 1}}
		err = methods.SetAssociations[models.Roles, models.Permissions](pgdb, &userRole, "Permissions", userRolePermissions)
		if err != nil {
			log.WithError(err).Fatal("При создании встроенных связей возникла не предвиденная ошибка:")
		}
	} else {
		log.Debugf("Объект \"%s\" уже существует в таблице \"roles\"", userRole.Name)
	}

	bAdmin := models.Users{Username: "admin", Password: "admin", Status: true}
	stat, err = methods.Create[models.Users](pgdb, map[string]interface{}{"username": "admin"}, &bAdmin)
	if err != nil {
		log.WithError(err).Fatal("При создании встроенного администратора возникла не предвиденная ошибка:")
	}
	if stat {
		log.Debugf("Объект \"%s\" успешно создан в таблице \"users\"", bAdmin.Username)
		bAdminRoles := []models.Roles{{ID: 1}, {ID: 2}}
		err = methods.SetAssociations[models.Users, models.Roles](pgdb, &bAdmin, "Roles", bAdminRoles)
		if err != nil {
			log.WithError(err).Fatal("При назначении ролей встроенному администратору возникла не предвиденная ошибка:")
		}
	} else {
		log.Debugf("Объект \"%s\" уже существует в таблице \"users\"", bAdmin.Username)
	}

	return nil
}
