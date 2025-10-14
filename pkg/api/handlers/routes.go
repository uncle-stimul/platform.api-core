package handlers

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	router.GET("/api/module/healthcheck", GetHealthcheck)
	router.GET("/api/module/info", GetModuleInfo)

	//router.GET("/api/auth/session", ###)
	//router.POST("/api/auth/session", ###)
	//router.DELETE("/api/auth/logout", ###)

	router.GET("/api/users/all", GetUsers)
	router.GET("/api/users/id", GetUser)
	router.POST("/api/users/id", CreateUser)
	router.PUT("/api/users/id", UpdateUser)
	router.DELETE("/api/users/id", DeleteUser)

	router.GET("/api/roles/all", selectRoles)
	router.GET("/api/roles/id", selectRole)
	router.POST("/api/roles/id", createRole)
	router.PUT("/api/roles/id", updateRole)
	router.DELETE("/api/roles/id", deleteRole)

	router.GET("/api/permissions/all", selectPermissions)
	router.GET("/api/permissions/id", selectPermission)
	router.POST("/api/permissions/id", createPermission)
	router.PUT("/api/permissions/id", updatePermission)
	router.DELETE("/api/permissions/id", deletePermission)

	router.GET("/api/sections/all", selectSections)
	router.GET("/api/sections/id", selectSection)
	router.POST("/api/sections/id", createSection)
	router.PUT("/api/sections/id", updateSection)
	router.DELETE("/api/sections/id", deleteSection)
}
