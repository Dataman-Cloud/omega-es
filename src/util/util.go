package util

import (
	"bytes"
	"crypto/md5"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Dataman-Cloud/omega-es/src/config"
	"github.com/Jeffail/gabs"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

const (
	EmailDefalutUser = "1"
	InternalTokenKey = "Sry-Svc-Token"
	LOG_ALARM_ID     = "Log-Alarm-Id"
)

func ReadBody(c *gin.Context) ([]byte, error) {
	b, err := ReadResponseBody(c.Request.Body)
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

func ReturnOKGin(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func ReturnCreatedOKGin(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

func ReturnOKObject(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}
func ReturnParamError(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, gin.H{"code": 17000, "data": "", "error": err})
}

func ReturnEsError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, gin.H{"code": 17001, "data": "", "error": err})
}

func ReturnDBError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, gin.H{"code": 17002, "data": "", "error": err})
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

func Header(c *gin.Context, key string) string {
	if values, _ := c.Request.Header[key]; len(values) > 0 {
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

func SendEmail(body string) error {
	req, err := http.NewRequest("POST", config.GetConfig().Murl, strings.NewReader(body))
	if err != nil {
		return err
	}
	token := CronTokenBuilder(EmailDefalutUser)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(InternalTokenKey, token)
	req.Header.Set("uid", EmailDefalutUser)
	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		return err
	}
	return nil
}

func AppScaling(body string, uid, clusterid, appid, alarmid int64) error {
	url := fmt.Sprintf("http://%s/api/v3/clusters/%d/apps/%d", config.GetConfig().Appurl, clusterid, appid)
	req, err := http.NewRequest("PATCH", url, strings.NewReader(body))
	if err != nil {
		return err
	}
	token := CronTokenBuilder(fmt.Sprintf("%d", uid))
	log.Debug("call app scaling:", url, body, token, uid)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(InternalTokenKey, token)
	req.Header.Set("Uid", fmt.Sprintf("%d", uid))
	req.Header.Set(LOG_ALARM_ID, fmt.Sprintf("%d", alarmid))
	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		return err
	}
	return nil
}

func DelScalingHistory(uid, alarmid int64) error {
	url := fmt.Sprintf("http://%s/api/v3/scales", config.GetConfig().Appurl)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	token := CronTokenBuilder(fmt.Sprintf("%d", uid))
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set(InternalTokenKey, token)
	req.Header.Set("Uid", fmt.Sprintf("%d", uid))
	req.Header.Set(LOG_ALARM_ID, fmt.Sprintf("%d", alarmid))
	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return err
	} else {
		body, err := ReadResponseBody(resp.Body)
		if err != nil {
			return err
		}
		jsonp, err := gabs.ParseJSON(body)
		if err != nil {
			return err
		}
		if code := jsonp.Path("code").Data().(float64); int64(code) != 0 {
			return errors.New(jsonp.Path("data").Data().(string))
		}
	}
	return nil
}
func GetInstance(uid, clusterid, appid int64) (int64, error) {
	url := fmt.Sprintf("http://%s/api/v3/clusters/%d/apps/%d", config.GetConfig().Appurl, clusterid, appid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	token := CronTokenBuilder(fmt.Sprintf("%d", uid))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(InternalTokenKey, token)
	req.Header.Set("Uid", fmt.Sprintf("%d", uid))
	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return 0, err
	} else {
		body, err := ReadResponseBody(resp.Body)
		if err != nil {
			return 0, err
		}
		jsonp, err := gabs.ParseJSON(body)
		if err != nil {
			return 0, err
		}
		if jsonp.Path("data.instances") == nil {
			return 0, nil
		}
		return int64(jsonp.Path("data.instances").Data().(float64)), nil

	}
	return 0, nil
}

func GetUserType(uid, clusterid int64) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", config.GetConfig().Userurl, clusterid), nil)
	log.Debugf("get usertype uri: %s", fmt.Sprintf("%s/%d", config.GetConfig().Userurl, clusterid))
	if err != nil {
		return "", err
	}
	token := CronTokenBuilder(fmt.Sprintf("%d", uid))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(InternalTokenKey, token)
	req.Header.Set("Uid", fmt.Sprintf("%d", uid))
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	respbody, err := ReadResponseBody(resp.Body)
	if err != nil {
		log.Errorf("get user tyep read reponse body error: %v", err)
		return "", err
	}
	jsonp, err := gabs.ParseJSON(respbody)
	if err != nil {
		return "", err
	}
	if jsonp.Path("data.group_id").Data() == nil {
		return "", errors.New("oweruser id")
	}
	return fmt.Sprintf("%d", int64(jsonp.Path("data.group_id").Data().(float64))), nil
}

func CronTokenBuilder(uid string) string {
	b64 := GetBaseEncoding()
	md5Token := fmt.Sprintf("%x", md5.Sum([]byte(uid)))
	b64Token := b64.EncodeToString([]byte(uid))
	token := b64.EncodeToString([]byte(fmt.Sprintf("%s:%s", md5Token, b64Token)))[:20]
	return strings.ToLower(token)
}

func GetBaseEncoding() *base64.Encoding {
	return base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
}

func ReplaceHtml(msg string) string {
	msg = strings.Replace(msg, "&", "&amp;", -1)
	msg = strings.Replace(msg, "<", "&lt;", -1)
	msg = strings.Replace(msg, ">", "&gt;", -1)
	msg = strings.Replace(msg, "\"", "&quot;", -1)
	msg = strings.Replace(msg, " ", "&nbsp;", -1)
	return msg
}
