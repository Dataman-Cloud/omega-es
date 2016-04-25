package util

import (
	"errors"
	"fmt"
	"github.com/Dataman-Cloud/omega-es/src/config"
	log "github.com/cihub/seelog"
	"github.com/jeffail/gabs"
	"net/http"
	"strings"
)

var url string

func init() {
	/*err, hosts := config.GetStringMapString("es", "hosts")
	if err != nil {
		log.Error(err)
	}
	err, port := config.GetStringMapString("es", "port")
	if err != nil {
		port = "9200"
		log.Warn("can't find es port default:9200")
	}*/
	//url = "http://" + strings.Split(hosts, ",")[0] + ":" + port + "/_watcher/watch/"
	url = fmt.Sprintf("http://%s:%d/_watcher/watch/", strings.Split(config.GetConfig().Ec.Hosts, ",")[0], config.GetConfig().Ec.Port)
}

func CrateWatcher(body, name string) error {
	req, err := http.NewRequest("PUT", url+name, strings.NewReader(body))
	if err != nil {
		log.Error("create watcher new request error: ", err)
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("create watcher client do error: ", err)
		return err
	}
	respbody, err := ReadResponseBody(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error("create watcher read response body error: ", err)
		return err
	}
	jsonp, err := gabs.ParseJSON(respbody)
	if err != nil {
		log.Error("create watcher response body parse to json error: ", err)
		return err
	}
	create, ok := jsonp.Path("created").Data().(bool)
	if !ok || !create {
		log.Error("can not create new watcher, may the name already exists")
		return errors.New("can not create new watcher, may the name already exists")
	}
	return nil
}

func GetWatcher(name string) (*gabs.Container, error) {
	resp, err := http.Get(url + name)
	if err != nil {
		log.Error("get watcher by name http error: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	respbody, err := ReadResponseBody(resp.Body)
	if err != nil {
		log.Error("get watcher by name read response body error: ", err)
		return nil, err
	}
	jsonp, err := gabs.ParseJSON(respbody)
	if err != nil {
		log.Error("get watcher by name response body parse json error: ", err)
		return nil, err
	}
	found, ok := jsonp.Path("found").Data().(bool)
	if !ok || !found {
		return nil, errors.New("can found watcher name: " + name)
	}
	return jsonp, err
}

func DeleteWatcherFromEs(name string) error {
	req, err := http.NewRequest("DELETE", url+name, nil)
	if err != nil {
		log.Error("delete watcher new request error: ", err)
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("delete watcher client do error: ", err)
		return nil
	}
	defer resp.Body.Close()
	respbody, err := ReadResponseBody(resp.Body)
	if err != nil {
		log.Error("delete watcher read response body error: ", err)
		return err
	}
	jsonp, err := gabs.ParseJSON(respbody)
	if err != nil {
		log.Error("delete watcher response body parse to json error: ", err)
		return err
	}
	found, ok := jsonp.Path("found").Data().(bool)
	if !ok {
		log.Error("delete watcher response body not found filed found")
		return errors.New("delete wathcer response body not found field found")
	}
	if !found {
		log.Error("delete watcher not found name: " + name)
		return errors.New("delete watcher not found name: " + name)
	}
	return nil
}
