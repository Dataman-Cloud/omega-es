package es

import (
	enjson "encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Dataman-Cloud/omega-es/src/cache"
	"github.com/Dataman-Cloud/omega-es/src/dao"
	"github.com/Dataman-Cloud/omega-es/src/model"
	. "github.com/Dataman-Cloud/omega-es/src/util"
	"github.com/Jeffail/gabs"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/sluu99/uuid"
)

const (
	EXPAND          = 1
	SHRINK          = 2
	STATUS_INACTIVE = 0
	STATUS_ACTIVE   = 1
	STATUS_UNUSABLE = 2
)

func CreateLogAlarm(c *gin.Context) {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("create log alarm can't get request body error: ", err)
		ReturnParamError(c, err.Error())
		return
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("create log alarm request parse to json error: ", err)
		ReturnParamError(c, err.Error())
		return
	}
	uid, ok := c.Get("uid")
	if !ok {
		log.Debug("create log alarm get uid error")
		ReturnParamError(c, "create log alarm get uid error")
		return
	}
	userid, err := strconv.ParseInt(uid.(string), 10, 64)
	if err != nil {
		log.Errorf("create log alarm parse userid to int64 error: %v", err)
		ReturnParamError(c, "create log alarm parse userid to int64 error: "+err.Error())
		return
	}
	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("create log alarm param can't get clusterid")
		ReturnParamError(c, "create log alarm param can't get clusterid")
		return
	}
	appid, ok := json.Path("appid").Data().(float64)
	if !ok {
		log.Error("create log alarm param can't get appid")
		ReturnParamError(c, "create log alarm param can't get appid")
		return
	}
	appalias, ok := json.Path("appalias").Data().(string)
	if !ok {
		log.Error("create log alarm param can't get appalias")
		ReturnParamError(c, "create log alarm param can't get appalias")
		return
	}
	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("create log alarm param can't get appname")
		ReturnParamError(c, "create log alarm param can't get appname")
		return
	}
	interval, ok := json.Path("interval").Data().(float64)
	if !ok {
		log.Error("create log alarm param can't get interval")
		ReturnParamError(c, "create log alarm param can't get interval")
		return
	}
	if int8(interval) <= 0 {
		log.Error("create log alarm interval must be greater than 0")
		ReturnParamError(c, "create log alarm interval must be greater than 0")
		return
	}

	gtnum, ok := json.Path("gtnum").Data().(float64)
	if !ok {
		log.Error("create log alarm param can't get gtnum")
		ReturnParamError(c, "create log alarm param can't get gtnum")
		return
	}

	if int8(gtnum) <= 0 {
		log.Error("create log alarm gtnum must be greater than 0")
		ReturnParamError(c, "create log alarm gtnum must be greater than 0")
		return
	}

	alarmname := uuid.Rand().Hex()
	usertype, ok := json.Path("usertype").Data().(string)
	if !ok {
		log.Error("create log alarm param can't get usertype")
		ReturnParamError(c, "create log alarm param can't get usertype")
		return
	}
	keyword, ok := json.Path("keyword").Data().(string)
	if !ok {
		log.Error("create log alarm param can't get keyword")
		ReturnParamError(c, "create log alarm param can't get keyword")
		return
	}
	emails, ok := json.Path("emails").Data().(string)
	if !ok {
		log.Error("create log alarm param can'get get emails")
		ReturnParamError(c, "create log alarm param can't get emails")
		return
	}

	scaling, ok := json.Path("scaling").Data().(bool)
	if !ok {
		log.Error("create log alarm param can't get scaling")
		ReturnParamError(c, "create log alarm param can't get scaling")
		return
	}
	maxs, ok := json.Path("maxs").Data().(float64)
	if !ok {
		maxs = 0
	}
	mins, ok := json.Path("mins").Data().(float64)
	if !ok {
		mins = 0
	}
	level, ok := json.Path("level").Data().(string)
	if !ok {
		level = ""
	}
	if count, err := dao.ExistAlarm(userid, usertype, alarmname); err != nil {
		log.Errorf("create log alarm judge alarm exist error: %v.", err)
		ReturnDBError(c, "cleate log alarm judge exist error: "+err.Error())
		return
	} else if count > 0 {
		log.Errorf("create log alarm %s already exist.", alarmname)
		ReturnParamError(c, "cleate log alarm already esist")
		return
	}

	if aexist, err := dao.GetAlarmByKeyword(int64(appid), userid, keyword, 0); err == nil {
		log.Debugf("get alarm by keyword count: %v", aexist)
		ReturnOKGin(c, map[string]interface{}{"code": 17018, "data": map[string]interface{}{"alarm": aexist}})
		return
	}

	//dao.AddAlarm
	aliasname := EncodAlias(alarmname, usertype, userid)
	alarm := &model.LogAlarm{
		Uid:        userid,
		Cid:        int64(clusterid),
		AppAlias:   appalias,
		AppName:    appname,
		Ival:       int64(interval),
		GtNum:      int64(gtnum),
		AlarmName:  alarmname,
		UserType:   usertype,
		KeyWord:    keyword,
		Emails:     emails,
		AliasName:  aliasname,
		CreateTime: time.Now(),
		Scaling:    scaling,
		Maxs:       uint64(maxs),
		Mins:       uint64(mins),
		AppId:      int64(appid),
		Level:      level,
	}

	if aid, err := dao.AddAlarm(alarm); err != nil {
		log.Errorf("create log alarm insert into alarm table error: %v", err)
		ReturnDBError(c, "create log alarm insert into alarm table error: "+err.Error())
		return
	} else {
		alarm.Id = aid
	}

	abody, err := enjson.Marshal(alarm)
	if err != nil {
		log.Errorf("create log alarm asearch parse to json string error: %v", err)
		ReturnDBError(c, "create log alarm asearch parse to json string error: "+err.Error())
		return
	}

	err = cache.AddAlarm(alarm.Id, abody)
	if err != nil {
		ReturnDBError(c, err.Error())
		return
	}

	ReturnCreatedOKGin(c, map[string]interface{}{"code": 0, "data": "create log alarm successful"})
	return
}

func GetAlarms(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("get alarms get param uid error")
		ReturnParamError(c, "get alarms get param uid error")
		return
	}
	userid, err := strconv.ParseInt(uid.(string), 10, 64)
	if err != nil {
		log.Errorf("get alarms parse userid to int64 error: %v", err)
		ReturnParamError(c, "get alarms parse userid to int64 error")
		return
	}
	pcount := c.Query("per_page")
	pagecount, _ := strconv.ParseInt(pcount, 10, 64)
	pnum := c.Query("page")
	pagenum, _ := strconv.ParseInt(pnum, 10, 64)
	order := c.Query("order")
	sortby := c.Query("sort_by")
	keyword := c.Query("keywords")

	alarms, err := dao.GetAlarmsByUser(userid, pagecount, pagenum, sortby, order, keyword)

	if err == nil {
		count, err := dao.CountAlarms(userid, keyword)
		if err != nil {
			log.Errorf("get alarms count error: %v", err)
			ReturnDBError(c, err.Error())
			return
		}
		ReturnOKGin(c, map[string]interface{}{"code": 0, "data": map[string]interface{}{"alarms": alarms, "count": count}})
		return
	}
	log.Errorf("get alarms error: %v", err)
	ReturnDBError(c, err.Error())
	return
}

func GetLogAlarm(c *gin.Context) {
	id := c.Param("id")
	jobid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Errorf("get alarm parse id to int64 error: %v", err)
		ReturnParamError(c, err.Error())
		return
	}
	alarm, err := dao.GetAlarmById(jobid)
	if err != nil {
		log.Errorf("get alarm by id error: %v", err)
		ReturnDBError(c, err.Error())
	}
	ReturnOKGin(c, map[string]interface{}{"code": 0, "data": map[string]interface{}{"alarm": alarm, "count": 1}})
}

func DeleteLogAlarm(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("delete alarm can't get uid")
		ReturnParamError(c, "delete alarm can't get uid")
		return
	}
	userid, err := strconv.ParseInt(uid.(string), 10, 64)
	if err != nil {
		log.Errorf("delete log alarm parse userid to int64 error: %v", err)
		ReturnParamError(c, err.Error())
		return
	}
	id := c.Param("id")
	jobid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Errorf("delete log alarm parse id to int64 error: %v", err)
		ReturnParamError(c, err.Error())
		return
	}
	alarm, err := dao.GetAlarmById(jobid)
	if err != nil {
		log.Errorf("delete log alarm can't get alarm by jobid")
		ReturnDBError(c, err.Error())
		return
	}

	if err = DelScalingHistory(userid, alarm.Id); err != nil {
		log.Errorf("delete scaling history error: %v", err)
		ReturnParamError(c, err.Error())
		return
	}
	if err = cache.DeleteAlarm(alarm.Id); err != nil {
		log.Errorf("delete cron error: %v", err)
		ReturnDBError(c, err.Error())
		return
	}

	if err = dao.DeleteAlarmByJobId(jobid); err != nil {
		log.Errorf("delete alarm error: %v", err)
		ReturnDBError(c, err.Error())
		return
	}
	ReturnOKObject(c, "delete alarm success")
	return
}

func JobExec(body []byte) error {
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("exec chronos job request parse to json error: ", err)
		return err
	}
	userid, ok := json.Path("uid").Data().(float64)
	if !ok {
		log.Error("exec chronos job param can't get userid")
		return errors.New("exec chronos jbo param can't get userid")
	}

	clusterid, ok := json.Path("cid").Data().(float64)
	if !ok {
		log.Error("exec chronos job param can't get clusterid")
		return errors.New("exec chronos job param can't get clusterid")
	}

	appalias, ok := json.Path("appalias").Data().(string)
	if !ok {
		log.Error("exec chronos job param can't get appname")
		return errors.New("exec chronos job param can't get appname")
	}

	keyword, ok := json.Path("keyword").Data().(string)
	if !ok {
		log.Error("exec chronos job param can't get keyword")
		return errors.New("exec chronos job param can't get keyword")
	}

	/*alarmname, ok := json.Path("alarmname").Data().(string)
	if !ok {
		log.Error("exec chronos job param can't get alarmname")
		return errors.New("exec chronos job param can't get alarmname")
	}*/

	interval, ok := json.Path("ival").Data().(float64)
	if !ok {
		log.Error("exec chronos job param can't get interval")
		return errors.New("exec chronos job param can't get interval")
	}

	/*usertype, ok := json.Path("usertype").Data().(string)
	if !ok {
		log.Error("exec chronos job param can't get usertype")
		return errors.New("exec chronos job param can't get usertype")
	}*/
	scaling, sok := json.Path("scaling").Data().(bool)
	if !ok {
		log.Error("exec chronos job param can't get scaling")
	}
	maxs, ok := json.Path("maxs").Data().(float64)
	if !ok {
		log.Error("exec chronos job param can't get maxs")
	}
	mins, ok := json.Path("mins").Data().(float64)
	if !ok {
		log.Error("exec chronos job param can't get mins")
	}

	appid, aok := json.Path("appid").Data().(float64)
	if !ok {
		log.Error("exec chronos job param can't get appid")
	}
	endtime := time.Now().Unix()
	starttime := endtime - int64(interval)*60
	query := `{"size":0,"query":{"filtered":{"query":{"bool":{"must":[{"term":{"clusterid":` + fmt.Sprintf("%d", int64(clusterid)) + `}},` +
		`{"term":{"typename":"` + appalias + `"}},{"match":{"msg":{"query":"` + keyword + `","analyzer":"ik"}}}]}},` +
		`"filter":{"bool":{"must":[{"range":{"timestamp":{"gte":"` + time.Unix(starttime, 0).Format(time.RFC3339) +
		`","lte":"` + time.Unix(endtime, 0).Format(time.RFC3339) + `"}}}]}}}},"aggs":{"ds":{"terms":{"field":"ipport","size":0}}}}`
	esindex := "dataman-app-" + fmt.Sprintf("%d", int64(clusterid)) + "-" + time.Now().String()[:10]
	estype := "dataman-" + appalias
	/*gid, err := GetUserType(int64(userid), int64(clusterid))
	if err == nil {
		esindex = "logstash-*" + gid + "-" + time.Now().String()[:10]
	}*/
	//out, err := Conn.Count(esindex, estype, nil, query)
	out, err := EsSearch(esindex, estype, nil, query)
	log.Debug("---", esindex, estype, query, err, string(out.RawJSON))
	if err != nil {
		log.Errorf("exec chronos job search es count error: %v", err)
		return err
	}
	//alarm, err := dao.GetAlarmByName(usertype, alarmname, int64(userid))
	alarm, err := dao.GetAlarmByName(int64(appid), int64(userid))
	if err != nil {
		log.Errorf("exec chronos job can't get alarm by alarmname error: %v", err)
		return err
	}
	if alarm.Isnotice != 1 {
		log.Debug("alarm isnotice is not equal to 1")
		return errors.New("alarm isnotice is not equal to 1")
	}
	rawjson, err := gabs.ParseJSON(out.RawJSON)
	if err != nil {
		log.Errorf("exec chronos job can't rawjson parse to json error: %v", err)
		return err
	}

	bs, err := rawjson.Path("aggregations.ds.buckets").Children()
	if err != nil {
		log.Errorf("exec chronos job get buckets children error: %v", err)
		return err
	}

	cache.UpdateScheduTime(alarm.Id)
	sore := SHRINK
	for _, b := range bs {
		if int64(b.Path("doc_count").Data().(float64)) >= alarm.GtNum {
			alarmHistory := &model.AlarmHistory{
				JobId:     alarm.Id,
				ExecTime:  time.Now(),
				ResultNum: int64(b.Path("doc_count").Data().(float64)),
				Uid:       int64(userid),
				Cid:       int64(clusterid),
				KeyWord:   keyword,
				AppName:   alarm.AppName,
				GtNum:     alarm.GtNum,
				Ival:      alarm.Ival,
				Ipport:    b.Path("key").Data().(string),
				Scaling:   alarm.Scaling,
				Maxs:      alarm.Maxs,
				Mins:      alarm.Mins,
				IsAlarm:   true,
				Level:     alarm.Level,
			}
			log.Debugf("---------: %s %s", alarm.AppName, alarm.Level)
			dao.AddAlaramHistory(alarmHistory)
			if sok && scaling {
				//shrinkorextend = true
				sore = EXPAND
			}
			break
		}
		if sok && scaling {
			//shrinkorextend = true
			sore = SHRINK
		}
	}
	log.Debug("-------:", alarm.AppName, "---", sok, scaling, aok, sore)
	if sok && scaling && aok {
		if sore == EXPAND {
			instances, err := GetInstance(int64(userid), int64(clusterid), int64(appid))
			log.Debug("++++++++++: ", instances, err, maxs)
			if err == nil && instances < int64(maxs) {
				sbody := gabs.New()
				sbody.Set("scale", "method")
				sbody.Set(uint64(maxs), "instances")
				err = AppScaling(sbody.String(), int64(userid), int64(clusterid), int64(appid), alarm.Id)
				log.Debugf("----add alarm scaling extend %v", err)
			}
		} else if sore == SHRINK {
			if instances, err := GetInstance(int64(userid), int64(clusterid), int64(appid)); err == nil && instances > int64(mins) && instances > 0 {
				sbody := gabs.New()
				sbody.Set("scale", "method")
				sbody.Set(instances-1, "instances")
				err = AppScaling(sbody.String(), int64(userid), int64(clusterid), int64(appid), alarm.Id)
				log.Debugf("----add alarm scaling shrink %v", err)
			}
		}
	}
	return nil
}

func GetAlarmHistory(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("get alarm history get param uid error")
		ReturnParamError(c, "get alarm history get param uid error")
		return
	}
	userid, err := strconv.ParseInt(uid.(string), 10, 64)
	if err != nil {
		log.Errorf("get alarm history parse userid to int64 error: %v", err)
		ReturnParamError(c, err.Error())
		return
	}
	id := c.Query("id")
	jobid, _ := strconv.ParseInt(id, 10, 64)
	_ = jobid
	pcount := c.Query("per_page")
	pagecount, _ := strconv.ParseInt(pcount, 10, 64)
	pnum := c.Query("page")
	pagenum, _ := strconv.ParseInt(pnum, 10, 64)
	order := c.Query("order")
	sortby := c.Query("sort_by")
	keyword := c.Query("keywords")
	historys, err := dao.GetHistoryByJobId(userid, pagecount, pagenum, sortby, order, keyword)
	if err != nil {
		log.Errorf("get alarm history error: %v", err)
		ReturnDBError(c, "get alarm history error: "+err.Error())
		return
	}
	count, err := dao.CountAlarmHistory(userid, keyword)
	if err != nil {
		log.Errorf("get alarm history error: %v", err)
		ReturnDBError(c, "get alarm history count error: "+err.Error())
		return
	}
	ReturnOKObject(c, map[string]interface{}{
		"events": historys,
		"count":  count,
	})
	return
}

func UpdateLogAlarm(c *gin.Context) {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("update log alarm can't get request body error: ", err)
		ReturnParamError(c, err.Error())
		return

	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("update log alarm request parse to json error: ", err)
		ReturnParamError(c, err.Error())
		return
	}

	id, ok := json.Path("id").Data().(float64)
	if !ok {
		log.Error("update log alarm param can't get id")
		ReturnParamError(c, "update log alarm param can't get id")
		return
	}

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("update log alarm param can't get clusterid")
		ReturnParamError(c, "update log alarm param can't get clusterid")
		return
	}

	appalias, ok := json.Path("appalias").Data().(string)
	if !ok {
		log.Error("update log alarm param can't get appalias")
		ReturnParamError(c, "update log alarm param can't get appalias")
		return
	}
	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("update log alarm param can't get appname")
		ReturnParamError(c, "update log alarm param can't get appname")
		return
	}
	interval, ok := json.Path("interval").Data().(float64)
	if !ok {
		log.Error("update log alarm param can't get interval")
		ReturnParamError(c, "update log alarm param can't get interval")
		return
	}
	if int64(interval) <= 0 {
		log.Error("update log alarm interval must be greater than 0")
		ReturnParamError(c, "update log alarm interval must be greater than 0")
		return
	}

	gtnum, ok := json.Path("gtnum").Data().(float64)
	if !ok {
		log.Error("update log alarm param can't get gtnum")
		ReturnParamError(c, "update log alarm param can't get gtnum")
		return
	}

	if int64(gtnum) <= 0 {
		log.Error("update log alarm gtnum must be greater than 0")
		ReturnParamError(c, "update log alarm gtnum must be greater than 0")
		return
	}

	usertype, ok := json.Path("usertype").Data().(string)
	if !ok {
		log.Error("update log alarm param can't get usertype")
		ReturnParamError(c, "update log alarm param can't get usertype")
		return
	}
	keyword, ok := json.Path("keyword").Data().(string)
	if !ok {
		log.Error("update log alarm param can't get keyword")
		ReturnParamError(c, "update log alarm param can't get keyword")
		return
	}
	emails, ok := json.Path("emails").Data().(string)
	if !ok {
		log.Error("update log alarm param can'get get emails")
		ReturnParamError(c, "update log alarm param can'get get emails")
		return
	}
	scaling, ok := json.Path("scaling").Data().(bool)
	if !ok {
		log.Error("update log alarm param can't get scaling")
		ReturnParamError(c, "update log alarm param can't get scaling")
		return
	}
	maxs, ok := json.Path("maxs").Data().(float64)
	if !ok {
		maxs = 0
	}
	mins, ok := json.Path("mins").Data().(float64)
	if !ok {
		mins = 0
	}
	level, ok := json.Path("level").Data().(string)
	if !ok {
		level = ""
	}

	alarm, err := dao.GetAlarmById(int64(id))
	if err != nil {
		log.Errorf("get alarm error: %v", err)
		ReturnDBError(c, err.Error())
		return
	}

	if aexist, err := dao.GetAlarmByKeyword(alarm.AppId, alarm.Uid, keyword, alarm.Id); err == nil {
		log.Debugf("get alarm by keyword count: %v", aexist)
		ReturnOKGin(c, map[string]interface{}{"code": 17018, "data": map[string]interface{}{"alarm": aexist}})
		return
	}
	alarm.Cid = int64(clusterid)
	alarm.AppAlias = appalias
	alarm.AppName = appname
	alarm.Ival = int64(interval)
	alarm.GtNum = int64(gtnum)
	alarm.UserType = usertype
	alarm.KeyWord = keyword
	alarm.Emails = emails
	alarm.Scaling = scaling
	alarm.Maxs = uint64(maxs)
	alarm.Mins = uint64(mins)
	alarm.Level = level
	if alarm.Isnotice == 1 {
		if err = cache.UpdateAlarm(&alarm); err != nil {
			log.Errorf("update alarm error")
			ReturnDBError(c, err.Error())
			return
		}
	}
	err = dao.UpdateAlarm(&alarm)
	if err != nil {
		log.Errorf("update alarm db table error: %v", err)
		ReturnDBError(c, err.Error())
		return
	}
	abody, err := enjson.Marshal(alarm)
	if err != nil {
		log.Errorf("update log alarm asearch parse to json string error: %v", err)
		ReturnDBError(c, err.Error())
		return
	}
	if err = cache.AddAlarm(alarm.Id, abody); err != nil {
		log.Errorf("restart alarm error")
		ReturnDBError(c, err.Error())
		return
	}

	ReturnOKObject(c, "update log alarm successfule")
	return
}

func StopLogAlarm(c *gin.Context) {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("operation alarm can't get request body error: ", err)
		ReturnParamError(c, err.Error())
		return
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("operation log alarm request parse to json error: ", err)
		ReturnParamError(c, err.Error())
		return
	}
	method, ok := json.Path("method").Data().(string)
	if !ok {
		log.Errorf("log alarm illegal request")
		ReturnParamError(c, "log alarm illegal request")
		return
	}
	id := c.Param("id")
	jobid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Errorf("stop log alarm parse id to int64 error: %v", err)
		ReturnParamError(c, "stop log alarm parse id to int64 error")
		return
	}
	alarm, err := dao.GetAlarmById(jobid)
	if err != nil {
		log.Errorf("stop log alarm get alarm by id error: %v", id)
		ReturnDBError(c, err.Error())
		return
	}
	if alarm.Isnotice == 2 {
		log.Debug("The service is not available")
		ReturnParamError(c, "The service is not available")
		return
	}
	if method == "stop" {
		if alarm.Isnotice == STATUS_INACTIVE {
			ReturnParamError(c, "alarm is already stop")
			return
		}
		alarm.Isnotice = STATUS_INACTIVE
		if err = cache.DeleteAlarm(alarm.Id); err != nil {
			log.Error("stop alarm error")
			ReturnDBError(c, err.Error())
			return
		}
		err = dao.UpdateAlarmStatus(&alarm)
		if err != nil {
			log.Errorf("stop log alarm update alarm status error: %v", err)
			ReturnDBError(c, err.Error())
			return
		}
	} else if method == "restart" {
		alarm.Isnotice = STATUS_ACTIVE
		err = dao.UpdateAlarmStatus(&alarm)
		if err != nil {
			log.Errorf("restart log alarm update alarm status error: %v", err)
			ReturnDBError(c, err.Error())
			return
		}
		abody, err := enjson.Marshal(alarm)
		if err != nil {
			log.Errorf("restart log alarm asearch parse to json string error: %v", err)
			ReturnDBError(c, err.Error())
			return
		}
		if err = cache.AddAlarm(alarm.Id, abody); err != nil {
			log.Errorf("restart alarm error")
			ReturnDBError(c, err.Error())
			return
		}

	} else {
		ReturnParamError(c, "illegality operation")
		return
	}
	ReturnOKObject(c, "operation alarm successful")
	return
}
