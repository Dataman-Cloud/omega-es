package es

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Dataman-Cloud/omega-es/src/util"
	"github.com/gin-gonic/gin"
	es "github.com/mattbaird/elastigo/lib"
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

var esindex string
var estype string

var fakeEsSearch = func(index string, _type string, args map[string]interface{}, query interface{}) (es.SearchResult, error) {
	esindex = index
	estype = _type
	return es.SearchResult{}, errors.New("FakeError")
}

func TestJobExec(t *testing.T) {
	util.EsSearch = fakeEsSearch
	var body = `{
						"uid":123,
						"cid":456,
						"appalias":"wwwww",
						"keyword":"key",
						"ival":1,
						"scaling":true,
						"maxs":10,
						"mins":1,
						"appid":10}
		`

	JobExec([]byte(body))
	assert.Equal(t, esindex, "dataman-app-456-"+time.Now().String()[:10])
	assert.Equal(t, estype, "dataman-wwwww")

}
