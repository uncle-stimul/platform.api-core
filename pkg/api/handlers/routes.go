package handlers

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	router.GET("/api/module/healthcheck", GetHealthcheck)
	router.GET("/api/module/info", GetModuleInfo)

	router.GET("/api/users/all", SelectUsers)
	router.GET("/api/users/id", SelectUser)
	router.POST("/api/users/id", CreateUser)
	router.PUT("/api/users/id", UpdateUser)
	router.DELETE("/api/users/id", DeleteUser)

	router.GET("/api/roles/all", SelectRoles)
	router.GET("/api/roles/id", SelectRole)
	router.POST("/api/roles/id", CreateRole)
	router.PUT("/api/roles/id", UpdateRole)
	router.DELETE("/api/roles/id", DeleteRole)

	router.GET("/api/permissions/all", SelectPermissions)
	router.GET("/api/permissions/id", SelectPermission)
	router.POST("/api/permissions/id", CreatePermission)
	router.PUT("/api/permissions/id", UpdatePermission)
	router.DELETE("/api/permissions/id", DeletePermission)

	router.GET("/api/sections/all", SelectSections)
	router.GET("/api/sections/id", SelectSection)
	router.POST("/api/sections/id", CreateSection)
	router.PUT("/api/sections/id", UpdateSection)
	router.DELETE("/api/sections/id", DeleteSection)
}
