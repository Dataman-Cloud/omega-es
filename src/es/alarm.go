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
	log "github.com/cihub/seelog"
	"github.com/jeffail/gabs"
	"github.com/labstack/echo"
	"github.com/sluu99/uuid"
	"strconv"
	"strings"
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
		Ival:       int8(interval),
		GtNum:      int64(gtnum),
		AlarmName:  alarmname,
		UserType:   usertype,
		KeyWord:    keyword,
		Emails:     emails,
		AliasName:  aliasname,
		CreateTime: time.Now(),
	}

	if aid, err := dao.AddAlarm(alarm); err != nil {
		log.Errorf("create log alarm insert into alarm table error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17006, "data": "create log alarm insert into alarm table error: " + err.Error()})
	} else {
		alarm.Id = aid
	}

	/*asearch := map[string]interface{}{
		"gtnum":     alarm.GtNum,
		"jobid":     alarm.Id,
		"userid":    userid,
		"clusterid": int64(clusterid),
		"keyword":   keyword,
		"appname":   appalias,
		"interval":  int8(interval),
		"usertype":  usertype,
		"alarmname": alarmname,
	}*/
	/*asearch := model.CronInfo{
		GtNum:     alarm.GtNum,
		JobId:     alarm.Id,
		UserId:    userid,
		ClusterId: int64(clusterid),
		KeyWord:   keyword,
		AppName:   appname,
		AppAlias:  appalias,
		Ival:      int8(interval),
		UserType:  usertype,
		AlarmName: alarmname,
	}*/

	abody, err := enjson.Marshal(alarm)
	if err != nil {
		log.Errorf("create log alarm asearch parse to json string error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17007, "data": "create log alarm asearch parse to json string error: " + err.Error()})
	}

	err = cache.AddAlarm(alarm.Id, abody)
	if err != nil {
		return ReturnError(c, map[string]interface{}{"code": 17006, "data": "add cron job error"})
	}
	/*authtoken := SchdulerAuth(usertype, alarmname, userid)
	scheduleCommand := fmt.Sprintf("curl -XPOST -H 'Content-Type: application/json' -H 'Authorization: %s' http://%s:%d/api/v3/scheduler -d '%s'", authtoken, config.GetConfig().Sh.Host, config.GetConfig().Sh.Port, string(abody))

	jobBody := map[string]interface{}{
		"name":     alarm.AliasName,
		"command":  scheduleCommand,
		"schedule": "R/" + time.Now().Format(time.RFC3339) + "/PT" + fmt.Sprintf("%d", int8(interval)) + "M",
		"owner":    "yqguo@dataman-inc.com",
		"async":    false,
	}
	//"schedule": "R/2016-04-24T09:03:31Z/PT" + fmt.Sprintf("%d", int8(interval)) + "S",
	cbody, _ := enjson.Marshal(jobBody)
	if err = CreateJob(string(cbody)); err != nil {
		log.Errorf("create log alarm add a chronos job error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17008, "data": "create log alarm add a chronos job error: " + err.Error()})
	}*/

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

//func JobExec(c *echo.Context) error {
func JobExec(body []byte) error {
	/*body, err := ReadBody(c)
	if err != nil {
		log.Error("exec chronos job cant't get request body")
		return ReturnError(c, map[string]interface{}{"code": 17001, "data": "exec chronos job can't get request body"})
	}*/
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("exec chronos job request parse to json error: ", err)
		return err
		//return ReturnError(c, map[string]interface{}{"code": 17002, "data": "exec chronos job request parse to json error"})
	}

	userid, ok := json.Path("uid").Data().(float64)
	if !ok {
		log.Error("exec chronos job param can't get userid")
		return errors.New("exec chronos jbo param can't get userid")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "exec chronos job param can't get userid"})
	}

	clusterid, ok := json.Path("cid").Data().(float64)
	if !ok {
		log.Error("exec chronos job param can't get clusterid")
		return errors.New("exec chronos job param can't get clusterid")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "exec chronos job param can't get clusterid"})
	}

	//appname, ok := json.Path("appname").Data().(string)
	appalias, ok := json.Path("appalias").Data().(string)
	if !ok {
		log.Error("exec chronos job param can't get appname")
		return errors.New("exec chronos job param can't get appname")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "exec chronos job param can't get appname"})
	}

	keyword, ok := json.Path("keyword").Data().(string)
	if !ok {
		log.Error("exec chronos job param can't get keyword")
		return errors.New("exec chronos job param can't get keyword")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "exec chronos job param can't get keyword"})
	}

	alarmname, ok := json.Path("alarmname").Data().(string)
	if !ok {
		log.Error("exec chronos job param can't get alarmname")
		return errors.New("exec chronos job param can't get alarmname")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "exec chronos job param can't get keyword"})
	}

	interval, ok := json.Path("ival").Data().(float64)
	if !ok {
		log.Error("exec chronos job param can't get interval")
		return errors.New("exec chronos job param can't get interval")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "exec chronos job param can't get inveral"})
	}

	usertype, ok := json.Path("usertype").Data().(string)
	if !ok {
		log.Error("exec chronos job param can't get usertype")
		return errors.New("exec chronos job param can't get usertype")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "exec chronos job param can't get usertype"})
	}

	/*authtoken := c.Request().Header.Get("Authorization")

	if authtoken != SchdulerAuth(usertype, alarmname, int64(userid)) {
		log.Error("Illegal request")
		//return ReturnError(c, map[string]interface{}{"code": 17003, "data": "Illegal request"})
	}*/

	endtime := time.Now().Unix()
	starttime := endtime - int64(interval)*60
	query := `{"query":{"bool":{"must":[{"match":{"msg":{"query":"` + keyword +
		`","analyzer":"ik"}}},{"range":{"timestamp":{"gte":"` + time.Unix(starttime, 0).Format(time.RFC3339) + `","lte":"` + time.Unix(endtime, 0).Format(time.RFC3339) + `"}}}]}}}`
	esindex := "logstash-*" + strconv.Itoa(int(userid)) + "-" + time.Now().String()[:10]
	gid, err := GetUserType(int64(userid), int64(clusterid))
	if err == nil {
		esindex = "logstash-*" + gid + "-" + time.Now().String()[:10]
	}
	estype := "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appalias
	out, err := Conn.Count(esindex, estype, nil, query)
	log.Debug("---", esindex, estype, query, out.Count)
	if err != nil {
		log.Errorf("exec chronos job search es count error: %v", err)
		return err
		//return ReturnError(c, map[string]interface{}{"code": 17010, "data": "exec chronos job search es count error: " + err.Error()})
	}
	alarm, err := dao.GetAlarmByName(usertype, alarmname, int64(userid))
	if err != nil {
		log.Errorf("exec chronos job can't get alarm by alarmname error: %v", err)
		return err
		//return ReturnError(c, map[string]interface{}{"code": 17011, "data": "exec chronos job can't get alarm by alarmname error: " + err.Error()})
	}
	cache.UpdateScheduTime(alarm.Id)
	alarmHistory := &model.AlarmHistory{
		JobId:     alarm.Id,
		ExecTime:  time.Now(),
		ResultNum: int64(out.Count),
		Uid:       int64(userid),
		Cid:       int64(clusterid),
		KeyWord:   keyword,
		AppName:   alarm.AppName,
		GtNum:     alarm.GtNum,
		Ival:      alarm.Ival,
	}
	if int64(out.Count) >= alarm.GtNum {
		alarmHistory.IsAlarm = true
		if aid, err := dao.AddAlaramHistory(alarmHistory); err != nil {
			log.Errorf("exec chronos job insert into alarm history error: %v", err)
			return errors.New("exec chronos job insert into alarm history error")
			//return ReturnError(c, map[string]interface{}{"code": 17012, "data": "exec chronos job insert into alarm history error: " + err.Error()})
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
	}
	return nil
	//return ReturnOK(c, map[string]interface{}{"code": 1, "data": "add alarm history successful"})
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
	if int8(interval) <= 0 {
		log.Error("update log alarm interval must be greater than 0")
		return ReturnError(c, map[string]interface{}{"code": 170015, "data": "update log alarm interval must be greater than 0"})
	}

	gtnum, ok := json.Path("gtnum").Data().(float64)
	if !ok {
		log.Error("update log alarm param can't get gtnum")
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "update log alarm param can't get gtnum"})
	}

	if int8(gtnum) <= 0 {
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
	alarm, err := dao.GetAlarmById(int64(id))
	if err != nil {
		log.Errorf("get alarm error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17016, "data": "get alarm error: " + err.Error()})
	}
	alarm.Cid = int64(clusterid)
	alarm.AppAlias = appalias
	alarm.AppName = appname
	alarm.Ival = int8(interval)
	alarm.GtNum = int64(gtnum)
	alarm.UserType = usertype
	alarm.KeyWord = keyword
	alarm.Emails = emails
	if err = cache.UpdateAlarm(&alarm); err != nil {
		log.Errorf("update alarm error")
		return ReturnError(c, map[string]interface{}{"code": 17015, "data": "update alarm error"})
	}
	err = dao.UpdateAlarm(&alarm)
	if err != nil {
		log.Errorf("update alarm db table error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17015, "data": "updata alarm db table error"})
	}
	/*asearch := map[string]interface{}{
		"gtnum":     alarm.GtNum,
		"jobid":     alarm.Id,
		"userid":    alarm.Uid,
		"clusterid": alarm.Cid,
		"keyword":   alarm.KeyWord,
		"appname":   alarm.AppAlias,
		"interval":  alarm.Ival,
		"usertype":  alarm.UserType,
		"alarmname": alarm.AlarmName,
	}*/
	/*asearch := model.CronInfo{
		GtNum:     alarm.GtNum,
		JobId:     alarm.Id,
		UserId:    alarm.Uid,
		ClusterId: alarm.Cid,
		KeyWord:   alarm.KeyWord,
		AppName:   alarm.AppName,
		AppAlias:  alarm.AppAlias,
		Ival:      alarm.Ival,
		UserType:  alarm.UserType,
		AlarmName: alarm.AlarmName,
	}*/
	abody, err := enjson.Marshal(alarm)
	if err != nil {
		log.Errorf("update log alarm asearch parse to json string error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17007, "data": "update log alarm asearch parse to json string error: " + err.Error()})
	}
	if err = cache.AddAlarm(alarm.Id, abody); err != nil {
		log.Errorf("restart alarm error")
		return ReturnError(c, map[string]interface{}{"code": 17016, "data": "restart alarm error"})
	}

	/*authtoken := SchdulerAuth(alarm.UserType, alarm.AlarmName, alarm.Uid)
	scheduleCommand := fmt.Sprintf("curl -XPOST -H 'Content-Type: application/json' -H 'Authorization: %s' http://%s:%d/api/v3/scheduler -d '%s'", authtoken, config.GetConfig().Sh.Host, config.GetConfig().Sh.Port, string(abody))

	jobBody := map[string]interface{}{
		"name":     alarm.AliasName,
		"command":  scheduleCommand,
		"schedule": "R/" + time.Now().Format(time.RFC3339) + "/PT" + fmt.Sprintf("%d", int8(interval)) + "M",
		"owner":    "yqguo@dataman-inc.com",
		"async":    false,
	}
	cbody, _ := enjson.Marshal(jobBody)*/
	/*authtoken := SchdulerAuth(usertype, alarmname, userid)
	_ = cbody
	/*if err = UpdateJob(string(cbody)); err != nil {
		log.Errorf("update log alarm add a chronos job error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17008, "data": "update log alarm add a chronos job error: " + err.Error()})
	}*/

	return ReturnOK(c, map[string]interface{}{"code": 0, "data": "update log alarm successful"})
}

func StopLogAlarm(c *echo.Context) error {
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
	/*err = DeleteJob(alarm.AliasName)
	if err != nil {
		log.Errorf("stop log alarm stop schduler job error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17003, "data": "stop log alarm stop schduler job error"})
	}*/
	return ReturnOK(c, map[string]interface{}{"code": 0, "data": "stop alarm successful"})
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
	/*if err = cache.AddAlarm(&alarm); err != nil {
		log.Errorf("restart alarm error")
		return ReturnError(c, map[string]interface{}{"code": 17016, "data": "restart alarm error"})
	}*/
	/*asearch := map[string]interface{}{
		"jobid":     alarm.Id,
		"gtnum":     alarm.GtNum,
		"userid":    alarm.Uid,
		"clusterid": alarm.Cid,
		"keyword":   alarm.KeyWord,
		"appname":   alarm.AppAlias,
		"interval":  alarm.Ival,
		"usertype":  alarm.UserType,
		"alarmname": alarm.AlarmName,
	}*/
	/*asearch := model.CronInfo{
		GtNum:     alarm.GtNum,
		JobId:     alarm.Id,
		UserId:    alarm.Uid,
		ClusterId: alarm.Cid,
		KeyWord:   alarm.KeyWord,
		AppName:   alarm.AppName,
		AppAlias:  alarm.AppAlias,
		Ival:      alarm.Ival,
		UserType:  alarm.UserType,
		AlarmName: alarm.AlarmName,
	}*/
	abody, err := enjson.Marshal(alarm)
	if err != nil {
		log.Errorf("restart log alarm asearch parse to json string error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17007, "data": "restart log alarm asearch parse to json string error: " + err.Error()})
	}
	if err = cache.AddAlarm(alarm.Id, abody); err != nil {
		log.Errorf("restart alarm error")
		return ReturnError(c, map[string]interface{}{"code": 17016, "data": "restart alarm error"})
	}
	/*authtoken := SchdulerAuth(alarm.UserType, alarm.AlarmName, alarm.Uid)
	scheduleCommand := fmt.Sprintf("curl -XPOST -H 'Content-Type: application/json' -H 'Authorization: %s' http://%s:%d/api/v3/scheduler -d '%s'", authtoken, config.GetConfig().Sh.Host, config.GetConfig().Sh.Port, string(abody))
	jobBody := map[string]interface{}{
		"name":     alarm.AliasName,
		"command":  scheduleCommand,
		"schedule": "R/" + time.Now().Format(time.RFC3339) + "/PT" + fmt.Sprintf("%d", alarm.Ival) + "M",
		"owner":    "yqguo@dataman-inc.com",
		"async":    false,
	}
	cbody, _ := enjson.Marshal(jobBody)
	_ = cbody*/
	/*if err = CreateJob(string(cbody)); err != nil {
		log.Errorf("restart log alarm add a chronos job error: %v", err)
		return ReturnError(c, map[string]interface{}{"code": 17008, "data": "restart log alarm add a chronos job error: " + err.Error()})
	}*/

	return ReturnOK(c, map[string]interface{}{"code": 0, "data": "restart log alarm successful"})
}
