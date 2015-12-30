package util

import (
	"bytes"
	"errors"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

func ReadBody(c *echo.Context) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request().Body)
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

func ReturnOK(ctx *echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, data)
}

func ReturnError(c *echo.Context, data interface{}) error {
	return c.JSON(http.StatusBadRequest, data)
}

func Header(c *echo.Context, key string) string {
	if values, _ := c.Request().Header[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}
