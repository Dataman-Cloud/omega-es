package es

import (
	//	"errors"
	"testing"

	"net/http"

	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	baseUrl string
	server  *httptest.Server
)

func TestMain(m *testing.M) {

	server := startHttpServer()
	baseUrl = server.URL
	defer server.Close()
	os.Exit(m.Run())
}

func auth(ctx *gin.Context) {
	ctx.Set("uid", 123)
	ctx.Next()
}

func startHttpServer() *httptest.Server {
	router := gin.New()
	v3 := router.Group("/api/v3", auth)
	{
		v3.POST("/es/index", SearchIndex)
	}
	return httptest.NewServer(router)
}

func TestSearchIndexBadRequest(t *testing.T) {

	req, err := http.NewRequest("POST", baseUrl+"/api/v3/es/index", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}
