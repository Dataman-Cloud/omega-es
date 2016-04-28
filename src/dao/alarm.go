package dao

import (
	"github.com/Dataman-Cloud/omega-es/src/model"
	"github.com/Dataman-Cloud/omega-es/src/util"
	log "github.com/cihub/seelog"
)

func ExistAlarm(uid int64, utype, alarmname string) (int, error) {
	db := util.DB()
	count := 0
	sql := `select count(*) from alarm where uid = ? and usertype = ? and alarmname = ?`
	err := db.Get(&count, sql, uid, utype, alarmname)
	return count, err
}

func CountAlarms(uid int64, keyword string) (int, error) {
	db := util.DB()

	count := 0
	sql := `select count(*) from alarm where uid = ? and appname like '%` + keyword + `%'`
	err := db.Get(&count, sql, uid)
	return count, err
}

func AddAlarm(alarm *model.LogAlarm) (int64, error) {
	db := util.DB()
	sql := `insert into alarm(uid, cid, appname, ival, gtnum, alarmname, usertype, keyword, emails, aliasname, createtime, appalias) values(:uid, :cid, :appname, :ival, :gtnum, :alarmname, :usertype, :keyword, :emails, :aliasname, :createtime, :appalias)`
	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		log.Error("insert into alarm error: ", err)
		return 0, err
	}
	defer func() {
		if stmt != nil {
			err = stmt.Close()
			if err != nil {
				log.Error("insert into alarm close stmt error: ", err)
			}
		}
	}()

	result, err := stmt.Exec(alarm)
	if err != nil {
		log.Error("insert into alarm error: ", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func GetAlarmsByUser(utype string, uid, pcount, pnum int64, sortby, order, keyword string) ([]model.LogAlarm, error) {
	db := util.DB()
	if pcount <= 0 {
		pcount = 10
	}
	if pnum <= 0 {
		pnum = 1
	}
	if sortby == "" {
		sortby = "createtime"
	}
	if order == "" {
		order = "asc"
	}

	alarms := []model.LogAlarm{}

	sql := `select * from alarm where uid = ?`
	if keyword != "" {
		sql = sql + ` and appname like '%` + keyword + `%'`
	}
	sql = sql + ` order by ` + sortby + ` ` + order + ` limit ?, ?`
	err := db.Select(&alarms, sql, uid, (pnum-1)*pcount, pcount)
	return alarms, err
}

func GetAlarmByName(utype, alarmname string, uid int64) (model.LogAlarm, error) {
	db := util.DB()
	alarm := model.LogAlarm{}
	sql := `select * from alarm where uid = ? and usertype = ? and alarmname = ?`
	err := db.Get(&alarm, sql, uid, utype, alarmname)
	return alarm, err
}

func GetAlarmById(id int64) (model.LogAlarm, error) {
	db := util.DB()
	alarm := model.LogAlarm{}
	sql := `select * from alarm where id = ?`
	err := db.Get(&alarm, sql, id)
	return alarm, err
}

func DeleteAlarmByJobId(jobid int64) error {
	db := util.DB()
	tx := db.MustBegin()
	_, err := tx.Exec("delete from alarm where id = ?", jobid)
	if err != nil {
		log.Errorf("delete alarm error: %v", err)
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("delete from alarmhistory where jobid = ?", jobid)
	if err != nil {
		log.Errorf("delete alarmhistor error: %v", err)
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Errorf("delete alarm commit error: %v", err)
		tx.Rollback()
		return err
	}
	return nil
}

func DeleteAlarmByUser(alarm *model.LogAlarm) error {
	db := util.DB()
	sql := `delete from alarm where uid = :uid and usertype = :usertype and alarmname = :alarmname`
	_, err := db.NamedExec(sql, alarm)
	return err
}

func AddAlaramHistory(ah *model.AlarmHistory) (int64, error) {
	db := util.DB()
	sql := `insert into alarmhistory(jobid, isalarm, exectime, resultnum) values (:jobid, :isalarm, :exectime, :resultnum)`
	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		log.Error("insert into alarm history error: ", err)
		return 0, err
	}
	defer func() {
		if stmt != nil {
			err = stmt.Close()
			if err != nil {
				log.Error("insert into alarm history close stmt error: ", err)
			}
		}
	}()

	result, err := stmt.Exec(ah)
	if err != nil {
		log.Error("insert into alarm history error: ", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func CountAlarmHistory(uid int64, keyword string) (int, error) {
	db := util.DB()

	count := 0
	sql := `select count(*) from alarmhistory as ah, alarm as a where ah.jobid = a.id and ah.isalarm = true and a.uid = ? and a.appname like '%` + keyword + `%'`
	err := db.Get(&count, sql, uid)
	return count, err
}

func GetHistoryByJobId(uid, pcount, pnum int64, sortby, order, keyword string) ([]model.AlarmHistory, error) {
	if pcount <= 0 {
		pcount = 10
	}
	if pnum <= 0 {
		pnum = 1
	}
	if sortby == "" || sortby == "exectime" {
		sortby = "ah.exectime"
	} else if sortby == "keyword" {
		sortby = "a.keyword"
	} else if sortby == "resultnum" {
		sortby = "ah.resultnum"
	} else if sortby == "jobid" {
		sortby = "ah.jobid"
	} else if sortby == "appname" {
		sortby = "a.appname"
	} else if sortby == "id" {
		sortby = "ah.id"
	}
	if order == "" {
		order = "asc"
	}
	db := util.DB()
	historys := []model.AlarmHistory{}
	sql := `select ah.id as id, ah.jobid as jobid, ah.isalarm as isalarm, ah.exectime as exectime, ah.resultnum as resultnum, a.appname as appname, a.gtnum as gtnum, a.ival as ival, a.keyword as keyword from alarmhistory as ah, alarm as a where ah.jobid = a.id and ah.isalarm = true and a.uid = ?`
	if keyword != "" {
		sql = sql + ` and a.appname like '%` + keyword + `%'`
	}
	sql = sql + ` order by ` + sortby + ` ` + order + ` limit ?,?`
	err := db.Select(&historys, sql, uid, (pnum-1)*pcount, pcount)
	return historys, err
}