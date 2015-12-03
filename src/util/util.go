package util

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ReadBody(c *gin.Context) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		return nil, errors.New("read request body error")
	}
	return buf.Bytes(), nil
}

func SameDay(start, end string) (bool, string) {
	s := strings.Split(start, "T")[0]
	e := strings.Split(end, "T")[0]
	if s == e {
		return true, s
	}
	return false, ""
}

func ReturnOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func ReturnError(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusBadRequest, data)
}

func Header(c *gin.Context, key string) string {
	if values, _ := c.Request.Header[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}
