package util

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fakeGetAppInstance(uid, clusterid, appid int64) (*http.Response, error) {
	res := &http.Response{}
	res.Body = &FakeBody{}
	return res, nil
}

type FakeBody struct {
}

func (fb *FakeBody) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (fb *FakeBody) Close() error {
	return nil
}

func TestGetInstance1(t *testing.T) {
	var fakeReadResponseBody1 = func(body io.ReadCloser) ([]byte, error) {
		return []byte(`{
			"data":{
					"instances":10
					}
				}
		`), nil
	}
	getAppInstance = fakeGetAppInstance
	ReadResponseBody = fakeReadResponseBody1
	number, _ := GetInstance(1, 11, 123)

	assert.Equal(t, number, int64(10))
}

func TestGetInstance2(t *testing.T) {
	var fakeReadResponseBody2 = func(body io.ReadCloser) ([]byte, error) {
		return []byte(`{
			"data":{
					"xxxxx":10
					}
				}
		`), nil
	}

	getAppInstance = fakeGetAppInstance
	ReadResponseBody = fakeReadResponseBody2
	number, _ := GetInstance(1, 11, 123)

	assert.Equal(t, number, int64(0))
}
