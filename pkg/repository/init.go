package repository

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"platform.api-core/pkg/models"
)

func initShema(pgdb *gorm.DB, log *logrus.Logger) {
	initTables(pgdb, log)
	initTablesEntities(pgdb, log)
	initTablesLinks(pgdb, log)

}

func initTables(pgdb *gorm.DB, log *logrus.Logger) {
	log.Debug("Инициализация SQL-схемы модуля api-core")
	if err := pgdb.AutoMigrate(
		&models.Users{},
		&models.Roles{},
		&models.Permissions{},
		&models.Sections{},
	); err != nil {
		log.WithError(err).Panicf("При инициализации SQL-схемы модуля api-core возникла критическая ошибка")
	}
}

func initTablesEntities(pgdb *gorm.DB, log *logrus.Logger) {
	log.Debug("Инициализация заполнения таблиц для модуля api-core")
	for _, section := range models.SchemaSections {
		if err := pgdb.Where("endpoint = ?", section.Endpoint).FirstOrCreate(&section).Error; err != nil {
			log.WithError(err).Errorf("Запись для объекта конечной точки \"%s\" уже существует", section.Endpoint)
		} else {
			log.Infof("Создана запись для объекта конечной точки \"%s\"", section.Endpoint)
		}
	}

	for _, permission := range models.SchemaPermissions {
		if err := pgdb.Where("name = ?", permission.Name).FirstOrCreate(&permission).Error; err != nil {
			log.WithError(err).Errorf("Запись для объекта прав доступа \"%s\" уже существует", permission.Name)
		} else {
			log.Infof("Создана запись для объекта прав доступа \"%s\"", permission.Name)
		}
	}

	for _, role := range models.SchemaRoles {
		if err := pgdb.Where("name = ?", role.Name).FirstOrCreate(&role).Error; err != nil {
			log.WithError(err).Errorf("Запись для объекта роли \"%s\" уже существует", role.Name)
		} else {
			log.Infof("Создана запись для объекта роли \"%s\"", role.Name)
		}
	}

	for _, user := range models.SchemaUsers {
		if err := pgdb.Where("Username = ?", user.Username).FirstOrCreate(&user).Error; err != nil {
			log.WithError(err).Errorf("Запись для объекта пользователя \"%s\" уже существует", user.Username)
		} else {
			log.Infof("Создана запись для объекта пользователя \"%s\"", user.Username)
		}
	}
}

func initTablesLinks(pgdb *gorm.DB, log *logrus.Logger) {
	for _, link := range models.SchemaLinks {
		log.Debugf("Инициализация заполнения связей между таблицами \"%s\" и \"%s\" для модуля api-core", link.ParentTable, link.ChildrenTable)
		var parentID uint
		queryParent := fmt.Sprintf("SELECT id FROM %s WHERE %s = ? LIMIT 1", link.ParentTable, link.ParentField)
		if err := pgdb.Raw(queryParent, link.ParentEntity).Scan(&parentID).Error; err != nil || parentID == 0 {
			log.WithError(err).Errorf("Объект-родитель \"%s\" не найден в таблице \"%s\", что привело к ошибке:", link.ParentEntity, link.ParentTable)
			continue
		}
		for _, childrenName := range link.ChildrenEntities {
			var childrenID uint
			queryChildren := fmt.Sprintf("SELECT id FROM %s WHERE %s = ? LIMIT 1", link.ChildrenTable, link.ChildrenField)
			if err := pgdb.Raw(queryChildren, childrenName).Scan(&childrenID).Error; err != nil || childrenID == 0 {
				log.WithError(err).Errorf("Объект-ребенок \"%s\" не найден в таблице \"%s\", что привело к ошибке:", childrenName, link.ChildrenTable)
				continue
			}

			linksTable := fmt.Sprintf("%s_%s", link.ParentTable, link.ChildrenTable)
			parentIDField := fmt.Sprintf("%s_id", link.ParentTable)
			childIDField := fmt.Sprintf("%s_id", link.ChildrenTable)

			query := fmt.Sprintf(
				"INSERT INTO %s (%s, %s) VALUES (?, ?) ON CONFLICT (%s, %s) DO NOTHING",
				linksTable, parentIDField, childIDField, parentIDField, childIDField,
			)

			if err := pgdb.Exec(query, parentID, childrenID).Error; err != nil {
				msg := fmt.Sprintf(
					"Ошибка при создании связи между объектами [ID: %d из таблицы %s] и [ID: %d из таблицы %s]",
					parentID, link.ParentTable, childrenID, link.ChildrenTable,
				)
				log.WithError(err).Error(msg)
			}
		}
	}
}
