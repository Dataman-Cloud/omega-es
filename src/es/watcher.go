package es

import (
	"encoding/base32"
	"github.com/Dataman-Cloud/omega-es/src/dao"
	"github.com/Dataman-Cloud/omega-es/src/model"
	. "github.com/Dataman-Cloud/omega-es/src/util"
	log "github.com/cihub/seelog"
	"github.com/jeffail/gabs"
	"github.com/labstack/echo"
	"strconv"
	"strings"
	"time"
)

func CreateWatcher(c *echo.Context) error {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("create watcher can't get request body")
		return ReturnError(c, map[string]interface{}{"code": 0, "error": "create watcher can't get request body"})
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("create watcher param parse json error")
		return ReturnError(c, map[string]interface{}{"code": 0, "error": "create watcher param parse json error"})
	}

	userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("create watcher param can't found userid")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher can't found userid"})
	}

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("create watcher param can't found clusterid")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher can't found clusterid"})
	}

	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("create watcher param can't found appname")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher can't found appname"})
	}

	watchername, ok := json.Path("watchername").Data().(string)
	if !ok {
		log.Error("create watcher param can't found watcher")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher can't found watcher"})
	}

	usertype, ok := json.Path("usertype").Data().(string)
	if !ok {
		log.Error("create watcher param can't found usertype")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher can't found usertype"})
	}

	interval, ok := json.Path("interval").Data().(string)
	if !ok {
		log.Error("create watcher param can't found interval")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher param can't found interval"})
	} else {
		_, err := strconv.Atoi(strings.Replace(interval, "m", "", 1))
		if err != nil {
			log.Error("wrong format monitoring cycle")
			return ReturnError(c, map[string]interface{}{"code": 0, "data": "wrong format monitoring cycle"})
		}
	}

	gtnum, ok := json.Path("gtnum").Data().(float64)
	if !ok {
		log.Error("create watcher param can't found gt number")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher param can't found gt number"})
	}

	timelimit, ok := json.Path("timelimit").Data().(string)
	if !ok {
		log.Error("create watcher param can't found timelimit")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher can't found timelimit"})
	}

	notifystr, ok := json.Path("notify").Data().(string)
	if !ok {
		log.Error("create watcher param can't found notify")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher can't found notify"})
	}
	notify, err := strconv.ParseBool(notifystr)
	if err != nil {
		log.Error("field notify not a bool type")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "field notfiy not a bool type"})
	}

	//log.Debug("this is watcher", userid, clusterid, appname, usertype, timelimit)
	count, err := dao.ExistWatcherByUidAndUtypeAndWname(int64(userid), usertype, watchername)
	if err != nil {
		log.Error("create check watcher exist error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "create check watcher exist error"})
	}
	/*if count > 0 {
		log.Debugf("create watcher already existing watcher name: %s", watchername)
		return ReturnOK(c, map[string]interface{}{"code": 0, "data": "watcher " + watchername + " already existing"})
	}*/
	if err == nil {
		esindex := "logstash-*" + strconv.Itoa(int(userid)) + "-" + time.Now().Format(time.RFC3339Nano)[:10]
		query := `{"query":{"bool":{"must":[{"term":{"clusterid":"` + strconv.Itoa(int(clusterid)) + `"}},{"term":{"typename":"` + appname + `"}},{"range":{"timestamp":{"gt":"now - ` + timelimit + `","lt":"now"}}}`
		if keyword, ok := json.Path("keyword").Data().(string); ok {
			query = query + `,{"match":{"msg":{"query":"` + keyword + `","analyzer":"ik"}}}`
		}
		if ipport, err := json.Path("ipport").Children(); err == nil && len(ipport) > 0 {
			var arr []string
			for _, ipp := range ipport {
				arr = append(arr, ipp.Data().(string))
			}
			query = query + `,{"terms":{"ipport":["` + strings.Join(arr, "\",\"") + `"],"minimum_match":1}}`
		}
		if source, err := json.Path("source").Children(); err == nil && len(source) > 0 {
			var arr []string
			for _, sour := range source {
				arr = append(arr, sour.Data().(string))
			}
			query = query + `,{"terms":{"source":["` + strings.Join(arr, "\",\"") + `"], "minimum_match":1}}`
		}
		query = query + `]}}}`
		//query = `{"query":{"match_all":{}}}`
		wbody := `{"trigger":{"schedule":{"interval":"` + interval + `"}},` +
			`"input":{"search":{"request":{"indices":["` + esindex + `"],` +
			`"body":` + query + `}}},` +
			`"condition":{"compare":{"ctx.payload.hits.total":{"gt":` + strconv.Itoa(int(gtnum)) + `}}},` +
			`"actions":{"my_webhook":{"webhook":{"method":"POST","host":"10.3.20.53","port":1323,"path":"/test/haha","headers": {"Content-type": "application/json"},` +
			`"body":"Encountered {{ctx.payload.hits.total}} errors"}}}}`
		watcher := &model.Watcher{
			Uid:    int64(userid),
			Utype:  usertype,
			Wname:  watchername,
			Wbody:  wbody,
			Cwname: base32.StdEncoding.EncodeToString([]byte(watchername + "_" + usertype + "_" + strconv.Itoa(int(userid)) + "_" + appname)),
			Notify: notify,
		}
		if count == 0 {
			_, err = dao.InsertWatcher(watcher)
			if err == nil && notify {
				//err = CrateWatcher(wbody, watchername+"_"+usertype+"_"+strconv.Itoa(int(userid)))
				err = CrateWatcher(wbody, watcher.Cwname)
				if err != nil {
					log.Error("create watcher http request error: ", err)
					return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher http request error"})
				}
			} else {
				log.Error("create watcher into mysql error: ", err)
				return ReturnError(c, map[string]interface{}{"code": 0, "data": "create watcher into mysql error"})
			}
		} else {
			err = dao.UpdateWatcher(watcher)
			if err == nil {
				//err = CrateWatcher(wbody, watchername+"_"+usertype+"_"+strconv.Itoa(int(userid)))
				if !notify {
					err = DeleteWatcherFromEs(watcher.Cwname)
					if err != nil {
						log.Error("update watcher http request error: ", err)
						return ReturnError(c, map[string]interface{}{"code": 0, "data": "update watcher http request error"})
					}
				}
				/*err = CrateWatcher(wbody, watcher.Cwname)
				if err != nil {
					log.Error("update watcher http request error: ", err)
					return ReturnError(c, map[string]interface{}{"code": 0, "data": "update watcher http request error"})
				}*/

			} else {
				log.Error("update watcher update mysql error: ", err)
				return ReturnError(c, map[string]interface{}{"code": 0, "data": "update watcher update mysql error"})
			}
		}
	}
	return ReturnOK(c, map[string]interface{}{"code": 1, "data": "create/or watcher successful"})
}

func GetWatchers(c *echo.Context) error {
	useridstr := c.Param("userid")
	usertype := c.Param("usertype")
	count := c.Param("count")
	pnum := c.Param("pnum")
	userid, err := strconv.ParseInt(useridstr, 10, 64)
	if err != nil {
		log.Error("get watcher userid string convert to int64 error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "get watchers userid string convert to int64 error"})
	}
	pagecount, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		log.Error("get watchers pagecount string convert to int64 error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "get watchers pagecount string convert to int64 error"})
	}
	pagenum, err := strconv.ParseInt(pnum, 10, 64)
	if err != nil {
		log.Error("get watcher pagenum string convert to int64 error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "get watchers pagenum string convert to int64 error"})
	}
	watchers, err := dao.GetWatchersByUser(userid, pagecount, pagenum, usertype)
	if err != nil {
		log.Error("get watchers list error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "get watchers list error"})
	}
	return ReturnOK(c, watchers)
}

func DeleteWatcher(c *echo.Context) error {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("delete watcher can't get request body")
		return ReturnError(c, map[string]interface{}{"code": 0, "error": "delete watcher can't get request body"})
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("delete watcher param parse json error")
		return ReturnError(c, map[string]interface{}{"code": 0, "error": "delete watcher param parse json error"})
	}

	userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("delete watcher param can't found userid")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "delete watcher can't found userid"})
	}
	watchername, ok := json.Path("watchername").Data().(string)
	if !ok {
		log.Error("delete watcher param can't found watchername")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "delete watcher can't found watchername"})
	}
	usertype, ok := json.Path("usertype").Data().(string)
	if !ok {
		log.Error("delete watcher param can't found usertype")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "delete watcher can't found usertype"})
	}
	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("delete watcher param can't found appname")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "delete watcher can't found appname"})
	}
	watcher := &model.Watcher{
		Uid:    int64(userid),
		Utype:  usertype,
		Wname:  watchername,
		Cwname: base32.StdEncoding.EncodeToString([]byte(watchername + "_" + usertype + "_" + strconv.Itoa(int(userid)) + "_" + appname)),
	}
	//err = DeleteWatcherFromEs(watchername + "_" + usertype + "_" + strconv.Itoa(int(userid)))
	err = DeleteWatcherFromEs(watcher.Cwname)
	if err != nil {
		//return ReturnError(c, map[string]interface{}{"code": 0, "data": err.Error()})
		log.Errorf("watcher not found watch name: %s", watcher.Cwname)
	}
	err = dao.DeleteWatcher(watcher)
	if err != nil {
		return ReturnError(c, map[string]interface{}{"code": 0, "data": err.Error()})
	}
	return ReturnOK(c, map[string]interface{}{"code": 0, "data": "delete watcher successful"})
}

func GetWatcherHistory(c *echo.Context) error {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("delete watcher can't get request body")
		return ReturnError(c, map[string]interface{}{"code": 0, "error": "delete watcher can't get request body"})
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("delete watcher param parse json error")
		return ReturnError(c, map[string]interface{}{"code": 0, "error": "delete watcher param parse json error"})
	}

	userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("delete watcher param can't found userid")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "delete watcher can't found userid"})
	}
	watchername, ok := json.Path("watchername").Data().(string)
	if !ok {
		log.Error("delete watcher param can't found watchername")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "delete watcher can't found watchername"})
	}
	usertype, ok := json.Path("usertype").Data().(string)
	if !ok {
		log.Error("delete watcher param can't found usertype")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "delete watcher can't found usertype"})
	}
	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("delete watcher param can't found appname")
		return ReturnError(c, map[string]interface{}{"code": 0, "data": "delete watcher can't found appname"})
	}
	cwname := base32.StdEncoding.EncodeToString([]byte(watchername + "_" + usertype + "_" + strconv.Itoa(int(userid)) + "_" + appname))
	esindex := ".watch_history*"
	query := `{"query":{"bool":{"must":[{"term":{"result.condition.met":true}},{"term":{"watch_id":"` + cwname + `"}}]}},"from":0,"size":10,
	"fields":["watch_id", "result.execution_time", "result.actions.webhook.request.body"],"sort":{"result.execution_time":"desc"}}`
	out, err := Conn.Search(esindex, "watch_record", nil, query)
	if err != nil {
		log.Error("searchindex search es error: ", err)
		return ReturnError(c, map[string]string{"error": "searchindex search es error"})
	}
	content, _ := gabs.ParseJSON(out.RawJSON)
	return ReturnOK(c, content.Data())
}
