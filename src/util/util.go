package util

import (
	"bytes"
	"errors"
	"github.com/labstack/echo"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/gommon/color"
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

func Logger() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()
			res := c.Response()
			logger := c.Echo().Logger()

			remoteAddr := req.RemoteAddr
			if ip := req.Header.Get(echo.XRealIP); ip != "" {
				remoteAddr = ip
			} else if ip = req.Header.Get(echo.XForwardedFor); ip != "" {
				remoteAddr = ip
			} else {
				remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
			}

			start := time.Now()
			if err := h(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			method := req.Method
			path := req.URL.Path
			if path == "" {
				path = "/"
			}
			size := res.Size()

			n := res.Status()
			code := color.Green(n)
			switch {
			case n >= 500:
				code = color.Red(n)
			case n >= 400:
				code = color.Yellow(n)
			case n >= 300:
				code = color.Cyan(n)
			}

			logger.Info("%s %s %s %s %s %s %d", start.Format(time.RFC3339Nano), remoteAddr, method, path, code, stop.Sub(start), size)
			return nil
		}
	}
}
