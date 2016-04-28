package util

import (
	"fmt"
	"github.com/Dataman-Cloud/omega-es/src/config"
	log "github.com/cihub/seelog"
	"net/http"
	"strings"
)

var chronosUrl string

const (
	addjob = "scheduler/iso8601"
	deljob = "scheduler/job/"
)

func init() {
	chronosUrl = fmt.Sprintf("http://%s:%d/", config.GetConfig().Ch.Host, config.GetConfig().Ch.Port)
}

func CreateJob(body string) error {
	req, err := http.NewRequest("POST", chronosUrl+addjob, strings.NewReader(body))
	if err != nil {
		log.Errorf("create chronos job new reques error: %v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Errorf("create chronos job client do error: %v", err)
	}
	/*respbody, err := ReadResponseBody(resp.Body)
	defer respbody.Body.Close()
	if err != nil {
		log.Errorf("create chronos job read response body error: %v", err)
		return err
	}

	jsonp, _ := gabs.ParseJSON(respbody)*/
	return nil
}

func DeleteJob(name string) error {
	req, err := http.NewRequest("DELETE", chronosUrl+deljob+name, nil)
	if err != nil {
		log.Errorf("delete chronos job error: %v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Errorf("delete chronos job do request error: %v", err)
		return err
	}
	return nil
}
