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

	router.GET("/api/roles/all", getRoles)
	//router.GET("/api/roles/id", getRole)
	//router.POST("/api/roles/id", postRole)
	//router.PUT("/api/roles/id", updRole)
	//router.DELETE("/api/roles/id", delRole)

	router.GET("/api/permissions/all", getPermissions)
	//router.GET("/api/permissions/id", getPermission)
	//router.POST("/api/permissions/id", postPermission)
	//router.PUT("/api/permissions/id", updPermission)
	//router.DELETE("/api/permissions/id", delPermission)

	//router.GET("/api/sections/all", getSections)
	//router.GET("/api/sections/id", getSection)
	//router.POST("/api/sections/id", postSection)
	//router.PUT("/api/sections/id", updSection)
	//router.DELETE("/api/sections/id", delSection)
}
