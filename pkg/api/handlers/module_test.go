package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"platform.api-core/pkg/api/handlers"
	"platform.api-core/pkg/models"
)

func TestHealthcheckHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/api/module/healthcheck", handlers.GetHealthcheck)

	req := httptest.NewRequest(http.MethodGet, "/api/module/healthcheck", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	re := regexp.MustCompile(`\s+`)
	checkRes := "{\"status\":\"running\"}"
	resBody := re.ReplaceAllString(w.Body.String(), "")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, resBody, checkRes)
}

func TestInfoHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/api/module/info", handlers.GetModuleInfo)

	req := httptest.NewRequest(http.MethodGet, "/api/module/info", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	re := regexp.MustCompile(`\s+`)
	checkRes := re.ReplaceAllString(
		fmt.Sprintf("{\"status\":\"success\",\"module\":\"%s\",\"version\":\"%s\",\"msg\":\"%s\"}",
			models.MetadataModule,
			models.MetadataVersion,
			models.MetadataMessage,
		), "",
	)
	resBody := re.ReplaceAllString(w.Body.String(), "")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, resBody, checkRes)
}
