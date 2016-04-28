package util

import (
	"bytes"
	"crypto/md5"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"strings"
)

func ReadBody(c *echo.Context) ([]byte, error) {
	b, err := ReadResponseBody(c.Request().Body)
	return b, err
}

func ReadResponseBody(body io.ReadCloser) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(body)
	if err != nil {
		return nil, errors.New("read request body error")
	}
	return buf.Bytes(), nil
}

func ReadBodyRequest(h *http.Request) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(h.Body)
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

func ReturnErrorResponse(w http.ResponseWriter, m map[string]interface{}) {
	b, _ := json.Marshal(m)
	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	w.Write(b)
}

func ReturnOKResponse(w http.ResponseWriter, m interface{}) {
	b, _ := json.Marshal(m)
	w.Header()["Content-Type"] = []string{"text/csv"}
	w.Write(b)
}

func Header(c *echo.Context, key string) string {
	if values, _ := c.Request().Header[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}

func FormatJson(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func EncodAlias(alarmname, usertype string, uid int64) string {
	alias := base32.StdEncoding.EncodeToString([]byte(alarmname + "_" + usertype + "_" + fmt.Sprintf("%d", uid)))
	return strings.Replace(strings.ToLower(alias), "=", "0", -1)
}

func SchdulerAuth(usertype, alarmname string, uid int64) string {
	an := fmt.Sprintf("%s-%s-%d", alarmname, usertype, uid)
	ana := fmt.Sprintf("%x", md5.Sum([]byte(an)))
	anb := base32.StdEncoding.EncodeToString([]byte(ana + "****"))
	anc := fmt.Sprintf("%x", md5.Sum([]byte(anb+"-+-+")))
	return anc

}
