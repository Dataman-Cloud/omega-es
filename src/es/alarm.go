package es

import (
	enjson "encoding/json"
	"errors"
	"fmt"
	"github.com/Dataman-Cloud/omega-es/src/cache"
	//"github.com/Dataman-Cloud/omega-es/src/config"
	"github.com/Dataman-Cloud/omega-es/src/dao"
	"github.com/Dataman-Cloud/omega-es/src/model"
	. "github.com/Dataman-Cloud/omega-es/src/util"
	"github.com/Jeffail/gabs"
	log "github.com/cihub/seelog"
	"github.com/labstack/echo"
	"github.com/sluu99/uuid"
	"strconv"
	//"strings"
	"time"
)

func CreateLogAlarm(c *echo.Context) error {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("create log alarm can't get request body error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 17001, "data": "create log alarm can't get request body"})
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("create log alarm request parse to json error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 17002, "data": "create log alarm request parse to json error"})
	}
	uid := c.Get("uid").(string)
	userid, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		log.Errorf("create log alarm parse userid to int64 error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm parse userid to int64 error: " + err.Error()})
	}
	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("create log alarm param can't get clusterid")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get clusterid"})
	}
	appid, ok := json.Path("appid").Data().(float64)
	if !ok {
		log.Error("create log alarm param can't get appid")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get appid"})
	}
	appalias, ok := json.Path("appalias").Data().(string)
	if !ok {
		log.Error("create log alarm param can't get appalias")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get appalias"})
	}
	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("create log alarm param can't get appname")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get appname"})
	}
	interval, ok := json.Path("interval").Data().(float64)
	if !ok {
		log.Error("create log alarm param can't get interval")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get interval"})
	}
	if int8(interval) <= 0 {
		log.Error("create log alarm interval must be greater than 0")
		return ReturnError(c, map[string]interface{}{"code": 170015, "data": "create log alarm interval must be greater than 0"})
	}

	gtnum, ok := json.Path("gtnum").Data().(float64)
	if !ok {
		log.Error("create log alarm param can't get gtnum")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get gtnum"})
	}

	if int8(gtnum) <= 0 {
		log.Error("create log alarm gtnum must be greater than 0")
		return ReturnError(c, map[string]interface{}{"code": 170015, "data": "create log alarm gtnum must be greater than 0"})
	}
	/*alarmname, ok := json.Path("alarmname").Data().(string)
	if !ok {
		log.Error("create log alarm param can't get alarm")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get alarmname"})
	}*/
	alarmname := uuid.Rand().Hex()
	usertype, ok := json.Path("usertype").Data().(string)
	if !ok {
		log.Error("create log alarm param can't get usertype")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get usertype"})
	}
	keyword, ok := json.Path("keyword").Data().(string)
	if !ok {
		log.Error("create log alarm param can't get keyword")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get keyword"})
	}
	emails, ok := json.Path("emails").Data().(string)
	if !ok {
		log.Error("create log alarm param can'get get emails")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get emails"})
	}

	scaling, ok := json.Path("scaling").Data().(bool)
	if !ok {
		log.Error("create log alarm param can't get scaling")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get scaling"})
	}
	maxs, ok := json.Path("maxs").Data().(float64)
	if !ok {
		maxs = 0
		//log.Error("create log alarm param can't get maxs")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get maxs"})
	}
	mins, ok := json.Path("mins").Data().(float64)
	if !ok {
		mins = 0
		//log.Error("create log alarm param can't get mins")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "create log alarm param can't get scaling"})
	}
	if count, err := dao.ExistAlarm(userid, usertype, alarmname); err != nil {
		log.Errorf("create log alarm judge alarm exist error: %v.", err)
		return ReturnError(c, map[string]interface{}{"code": 17004, "data": "cleate log alarm judge exist error: " + err.Error()})
	} else if count > 0 {
		log.Errorf("create log alarm %s already exist.", alarmname)
		return ReturnError(c, map[string]interface{}{"code": 17005, "data": "cleate log alarm already esist"})
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
	}

	if aid, err := dao.AddAlarm(alarm); err != nil {
		log.Errorf("create log alarm insert into alarm table error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17006, "data": "create log alarm insert into alarm table error: " + err.Error()})
	} else {
		alarm.Id = aid
	}

	abody, err := enjson.Marshal(alarm)
	if err != nil {
		log.Errorf("create log alarm asearch parse to json string error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17007, "data": "create log alarm asearch parse to json string error: " + err.Error()})
	}

	err = cache.AddAlarm(alarm.Id, abody)
	if err != nil {
		return ReturnError(c, map[string]interface{}{"code": 17006, "data": "add cron job error"})
	}

	return ReturnOK(c, map[string]interface{}{"code": 0, "data": "create log alarm successful"})
}

func GetAlarms(c *echo.Context) error {
	//utype := c.Query("usertype")
	//uid := c.Query("uid")
	utype := ""
	//userid, _ := strconv.ParseInt(uid, 10, 64)
	uid := c.Get("uid").(string)
	userid, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		log.Errorf("get alarms parse userid to int64 error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "get alarms parse userid to int64 error"})
	}
	pcount := c.Query("per_page")
	pagecount, _ := strconv.ParseInt(pcount, 10, 64)
	pnum := c.Query("page")
	pagenum, _ := strconv.ParseInt(pnum, 10, 64)
	order := c.Query("order")
	sortby := c.Query("sort_by")
	keyword := c.Query("keywords")

	alarms, err := dao.GetAlarmsByUser(utype, userid, pagecount, pagenum, sortby, order, keyword)

	if err == nil {
		count, err := dao.CountAlarms(userid, keyword)
		if err != nil {
			log.Error("get alarms count error")
			return ReturnError(c, map[string]interface{}{"code": 17003, "data": "get alarms count error"})
		}
		return ReturnOK(c, map[string]interface{}{"code": 0, "data": map[string]interface{}{"alarms": alarms, "count": count}})
	}
	log.Errorf("get alarms error: %v", err)
	return ReturnError(c, map[string]interface{}{"code": 17009, "data": "get alarms error: " + err.Error()})
}

func GetLogAlarm(c *echo.Context) error {
	id := c.Param("id")
	jobid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Errorf("get alarm parse id to int64 error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "get alarm parse id to int64 error"})
	}
	alarm, err := dao.GetAlarmById(jobid)
	if err != nil {
		log.Errorf("get alarm by id error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "get alarm by id error"})
	}
	return ReturnOK(c, map[string]interface{}{"code": 0, "data": map[string]interface{}{"alarm": alarm, "count": 1}})
}

func DeleteLogAlarm(c *echo.Context) error {
	id := c.Param("id")
	jobid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Errorf("delete log alarm parse id to int64 error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "delete log alarm parse id to int64 error:: " + err.Error()})
	}
	alarm, err := dao.GetAlarmById(jobid)
	if err != nil {
		log.Errorf("delete log alarm can't get alarm by jobid")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "delete log alarm can't get alarm by jobid error: " + err.Error()})
	}

	/*aliasname := EncodAlias(alarm.AlarmName, alarm.UserType, alarm.Uid)
	if err = DeleteJob(aliasname); err != nil {
		log.Errorf("delete alarm remove chronos job error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "delete alarm remove chronos job error: " + err.Error()})
	}*/

	if err = cache.DeleteAlarm(alarm.Id); err != nil {
		log.Errorf("delete cron error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "delete cron job error"})
	}

	if err = dao.DeleteAlarmByJobId(jobid); err != nil {
		log.Errorf("delete alarm error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "delete alarm error: " + err.Error()})
	}
	return ReturnOK(c, map[string]interface{}{"code": 0, "data": "delete alarm successful"})
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

	alarmname, ok := json.Path("alarmname").Data().(string)
	if !ok {
		log.Error("exec chronos job param can't get alarmname")
		return errors.New("exec chronos job param can't get alarmname")
	}

	interval, ok := json.Path("ival").Data().(float64)
	if !ok {
		log.Error("exec chronos job param can't get interval")
		return errors.New("exec chronos job param can't get interval")
	}

	usertype, ok := json.Path("usertype").Data().(string)
	if !ok {
		log.Error("exec chronos job param can't get usertype")
		return errors.New("exec chronos job param can't get usertype")
	}
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
	esindex := "logstash-*" + strconv.Itoa(int(userid)) + "-" + time.Now().String()[:10]
	/*gid, err := GetUserType(int64(userid), int64(clusterid))
	if err == nil {
		esindex = "logstash-*" + gid + "-" + time.Now().String()[:10]
	}*/
	estype := "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appalias
	esindex = "*"
	estype = ""
	//out, err := Conn.Count(esindex, estype, nil, query)
	out, err := Conn.Search(esindex, estype, nil, query)
	log.Debug("---", esindex, estype, query, err, string(out.RawJSON))
	if err != nil {
		log.Errorf("exec chronos job search es count error: %v", err)
		return err
	}
	alarm, err := dao.GetAlarmByName(usertype, alarmname, int64(userid))
	if err != nil {
		log.Errorf("exec chronos job can't get alarm by alarmname error: %v", err)
		return err
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
	if len(bs) == 0 {
		return nil
	}
	shrinkorextend := false
	sore := 0
	for _, b := range bs {
		if int64(b.Path("doc_count").Data().(float64)) < alarm.GtNum {
			if sok && scaling {
				shrinkorextend = true
				sore = 2
			}
			break
		}
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
		}
		dao.AddAlaramHistory(alarmHistory)
		if sok && scaling {
			shrinkorextend = true
			sore = 1
		}
	}
	log.Debug("-------:", alarm.AppName, "---", shrinkorextend, aok)
	if shrinkorextend && aok {
		if sore == 1 {
			instances, err := GetInstance(int64(userid), int64(clusterid), int64(appid))
			if err == nil && instances != int64(maxs) {
				sbody := gabs.New()
				sbody.Set("scale", "method")
				sbody.Set(uint64(maxs), "instances")
				err = AppScaling(sbody.String(), int64(userid), int64(clusterid), int64(appid), alarm.Id)
				log.Debugf("----add alarm scaling extend %v", err)
			}
		} else if sore == 2 {
			if instances, err := GetInstance(int64(userid), int64(clusterid), int64(appid)); err == nil && instances > int64(mins) && instances > 0 {
				sbody := gabs.New()
				sbody.Set("scale", "method")
				sbody.Set(instances-1, "instances")
				err = AppScaling(sbody.String(), int64(userid), int64(clusterid), int64(appid), alarm.Id)
				log.Debugf("----add alarm scaling shrink %v", err)
			}
		}
	}
	/*alarmHistory := &model.AlarmHistory{
		JobId:     alarm.Id,
		ExecTime:  time.Now(),
		ResultNum: int64(out.Count),
		Uid:       int64(userid),
		Cid:       int64(clusterid),
		KeyWord:   keyword,
		AppName:   alarm.AppName,
		GtNum:     alarm.GtNum,
		Ival:      alarm.Ival,
		Ipport:    alarm.Ipport,
		Scaling:   alarm.Scaling,
		Maxs:      alarm.Maxs,
		Mins:      alarm.Mins,
	}
	if int64(out.Count) >= alarm.GtNum {
		alarmHistory.IsAlarm = true
		if aid, err := dao.AddAlaramHistory(alarmHistory); err != nil {
			log.Errorf("exec chronos job insert into alarm history error: %v", err)
			return errors.New("exec chronos job insert into alarm history error")
		} else {
			memail := map[string]interface{}{
				"template": "alarm",
				"subject":  fmt.Sprintf("数人云日志告警-策略%d-告警事件%d", alarm.Id, aid),
				"emails":   strings.Split(alarm.Emails, ","),
				"data": map[string]string{
					"content": fmt.Sprintf("应用%s日志在%d分钟内出现关键词%s%d次，请您关注", alarm.AppName, alarm.Ival, alarm.KeyWord, int64(out.Count)),
				},
			}
			bemail, err := enjson.Marshal(memail)
			if err == nil {
				err = SendEmail(string(bemail))
				if err != nil {
					log.Errorf("alarm send email error: %v", err)
					return errors.New("alarm send email error")
				}
			}
		}
	} else {
		alarmHistory.IsAlarm = false
	}*/
	return nil
}

func GetAlarmHistory(c *echo.Context) error {
	uid := c.Get("uid").(string)
	userid, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		log.Errorf("get alarm history parse userid to int64 error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "get alarm history parse userid to int64 error: " + err.Error()})
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
		return ReturnError(c, map[string]interface{}{"code": 17013, "data": "get alarm history error: " + err.Error()})
	}
	count, err := dao.CountAlarmHistory(userid, keyword)
	if err != nil {
		log.Errorf("get alarm history error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17014, "data": "get alarm history count error"})
	}
	return ReturnOK(c, map[string]interface{}{"code": 0, "data": map[string]interface{}{
		"events": historys,
		"count":  count,
	}})
}

func UpdateLogAlarm(c *echo.Context) error {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("update log alarm can't get request body error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 17001, "data": "update log alarm can't get request body"})
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("update log alarm request parse to json error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 17002, "data": "update log alarm request parse to json error"})
	}

	id, ok := json.Path("id").Data().(float64)
	if !ok {
		log.Error("update log alarm param can't get id")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get id"})
	}

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("update log alarm param can't get clusterid")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get clusterid"})
	}

	appalias, ok := json.Path("appalias").Data().(string)
	if !ok {
		log.Error("update log alarm param can't get appalias")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get appalias"})
	}
	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("update log alarm param can't get appname")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get appname"})
	}
	interval, ok := json.Path("interval").Data().(float64)
	if !ok {
		log.Error("update log alarm param can't get interval")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get interval"})
	}
	if int64(interval) <= 0 {
		log.Error("update log alarm interval must be greater than 0")
		return ReturnError(c, map[string]interface{}{"code": 170015, "data": "update log alarm interval must be greater than 0"})
	}

	gtnum, ok := json.Path("gtnum").Data().(float64)
	if !ok {
		log.Error("update log alarm param can't get gtnum")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get gtnum"})
	}

	if int64(gtnum) <= 0 {
		log.Error("update log alarm gtnum must be greater than 0")
		return ReturnError(c, map[string]interface{}{"code": 170015, "data": "update log alarm gtnum must be greater than 0"})
	}

	usertype, ok := json.Path("usertype").Data().(string)
	if !ok {
		log.Error("update log alarm param can't get usertype")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get usertype"})
	}
	keyword, ok := json.Path("keyword").Data().(string)
	if !ok {
		log.Error("update log alarm param can't get keyword")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get keyword"})
	}
	emails, ok := json.Path("emails").Data().(string)
	if !ok {
		log.Error("update log alarm param can'get get emails")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get emails"})
	}
	scaling, ok := json.Path("scaling").Data().(bool)
	if !ok {
		log.Error("update log alarm param can't get scaling")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get scaling"})
	}
	maxs, ok := json.Path("maxs").Data().(float64)
	if !ok {
		maxs = 0
		//log.Error("update log alarm param can't get maxs")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get maxs"})
	}
	mins, ok := json.Path("mins").Data().(float64)
	if !ok {
		mins = 0
		//log.Error("update log alarm param can't get mins")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get scaling"})
	}
	alarm, err := dao.GetAlarmById(int64(id))
	if err != nil {
		log.Errorf("get alarm error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17016, "data": "get alarm error: " + err.Error()})
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
	if err = cache.UpdateAlarm(&alarm); err != nil {
		log.Errorf("update alarm error")
		return ReturnError(c, map[string]interface{}{"code": 17015, "data": "update alarm error"})
	}
	err = dao.UpdateAlarm(&alarm)
	if err != nil {
		log.Errorf("update alarm db table error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17015, "data": "updata alarm db table error"})
	}
	abody, err := enjson.Marshal(alarm)
	if err != nil {
		log.Errorf("update log alarm asearch parse to json string error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17007, "data": "update log alarm asearch parse to json string error: " + err.Error()})
	}
	if err = cache.AddAlarm(alarm.Id, abody); err != nil {
		log.Errorf("restart alarm error")
		return ReturnError(c, map[string]interface{}{"code": 17016, "data": "restart alarm error"})
	}

	return ReturnOK(c, map[string]interface{}{"code": 0, "data": "update log alarm successful"})
}

func StopLogAlarm(c *echo.Context) error {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("operation alarm can't get request body error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 17001, "data": "operation log alarm can't get request body"})
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("operation log alarm request parse to json error: ", err)
		return ReturnError(c, map[string]interface{}{"code": 17002, "data": "operation log alarm request parse to json error"})
	}
	method, ok := json.Path("method").Data().(string)
	if !ok {
		log.Errorf("log alarm illegal request")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "log alarm illegal request"})
	}
	id := c.Param("id")
	jobid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Errorf("stop log alarm parse id to int64 error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "stop log alarm parse id to int64 error"})
	}
	alarm, err := dao.GetAlarmById(jobid)
	if err != nil {
		log.Errorf("stop log alarm get alarm by id error: %v", id)
		return ReturnError(c, map[string]interface{}{"code": 17016, "data": "stop log alarm get alarm by id error"})
	}
	if method == "stop" {
		if !alarm.Isnotice {
			return ReturnOK(c, map[string]interface{}{"code": 0, "data": "alarm is already stop"})
		}
		alarm.Isnotice = false
		if err = cache.DeleteAlarm(alarm.Id); err != nil {
			log.Error("stop alarm error")
			return ReturnError(c, map[string]interface{}{"code": 17016, "data": "stop alarm error"})
		}
		err = dao.UpdateAlarmStatus(&alarm)
		if err != nil {
			log.Errorf("stop log alarm update alarm status error: %v", err)
			return ReturnError(c, map[string]interface{}{"code": 17016, "data": "stop log alarm update alarm status error"})
		}
		//return ReturnOK(c, map[string]interface{}{"code": 0, "data": "stop alarm successful"})
	} else if method == "restart" {
		alarm.Isnotice = true
		err = dao.UpdateAlarmStatus(&alarm)
		if err != nil {
			log.Errorf("restart log alarm update alarm status error: %v", err)
			return ReturnError(c, map[string]interface{}{"code": 17016, "data": "restart log alarm update alarm status error"})
		}
		abody, err := enjson.Marshal(alarm)
		if err != nil {
			log.Errorf("restart log alarm asearch parse to json string error: %v", err)
			return ReturnError(c, map[string]interface{}{"code": 17007, "data": "restart log alarm asearch parse to json string error: " + err.Error()})
		}
		if err = cache.AddAlarm(alarm.Id, abody); err != nil {
			log.Errorf("restart alarm error")
			return ReturnError(c, map[string]interface{}{"code": 17016, "data": "restart alarm error"})
		}

		//return ReturnOK(c, map[string]interface{}{"code": 0, "data": "restart log alarm successful"})
	} else {
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "illegality operation"})
	}
	return ReturnOK(c, map[string]interface{}{"code": 0, "data": "operation alarm successful"})
}

func RestartLogAlarm(c *echo.Context) error {
	id := c.Param("id")
	jobid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Errorf("restart log alarm parse id to int64 error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "restart log alarm parse id to int64 error"})
	}
	alarm, err := dao.GetAlarmById(jobid)
	if err != nil {
		log.Errorf("restart log alarm can't get alarm by jobid")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "restart log alarm can't get alarm by jobid error: " + err.Error()})
	}
	alarm.Isnotice = true
	err = dao.UpdateAlarmStatus(&alarm)
	if err != nil {
		log.Errorf("restart log alarm update alarm status error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17016, "data": "restart log alarm update alarm status error"})
	}
	abody, err := enjson.Marshal(alarm)
	if err != nil {
		log.Errorf("restart log alarm asearch parse to json string error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17007, "data": "restart log alarm asearch parse to json string error: " + err.Error()})
	}
	if err = cache.AddAlarm(alarm.Id, abody); err != nil {
		log.Errorf("restart alarm error")
		return ReturnError(c, map[string]interface{}{"code": 17016, "data": "restart alarm error"})
	}

	return ReturnOK(c, map[string]interface{}{"code": 0, "data": "restart log alarm successful"})
}
