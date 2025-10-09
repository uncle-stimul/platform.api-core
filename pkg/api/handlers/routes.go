package handlers

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	router.GET("/healthcheck", getHealthcheck)
	router.GET("/api/status", getStatus)

	//router.POST("/api/auth/login", ###)
	//router.POST("/api/auth/logout", ###)

	router.GET("/api/users/all", getUsers)
	//router.GET("/api/users/id", getUser)
	//router.POST("/api/users/id", postUser)
	//router.PUT("/api/users/id", updUser)
	//router.DELETE("/api/users/id", delUser)

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
