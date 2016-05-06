package util

import (
	"fmt"
	"github.com/Dataman-Cloud/omega-es/src/config"
	"github.com/Dataman-Cloud/omega-es/src/model"
	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattes/migrate/driver/mysql"
	"github.com/mattes/migrate/migrate"
	"sync"
)

func init() {
	DB()
	upgradeDB()
}

var db *sqlx.DB

func DB() *sqlx.DB {
	if db != nil {
		return db
	}
	mutex := sync.Mutex{}
	mutex.Lock()
	db, _ = InitDB()
	defer mutex.Unlock()
	return db
}

func upgradeDB() {
	uri := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		config.GetConfig().Mc.UserName,
		config.GetConfig().Mc.PassWord,
		config.GetConfig().Mc.Host,
		config.GetConfig().Mc.Port,
		config.GetConfig().Mc.DataBase)
	//uri = "mysql://root:111111@tcp(10.3.37.6:3306)/oapp?parseTime=true&loc=Local"
	log.Info("upgrade db mysql drive: ", uri)
	errors, ok := migrate.UpSync(uri, "./sql")
	if errors != nil && len(errors) > 0 {
		for _, err := range errors {
			log.Error("db err", err)
		}
		log.Error("can't upgrade db", errors)
		log.Flush()
		//panic(-1)
	}
	if !ok {
		log.Error("can't upgrade db")
		log.Flush()
		//panic(-1)
	}
	log.Info("DB upgraded")
	log.Flush()
	updatedata()
}

func updatedata() {
	log.Debug("update start...")
	historys := []model.AlarmHistory{}
	sql := `select ah.id as id, ah.jobid as jobid, ah.isalarm as isalarm, ah.exectime as exectime, ah.resultnum as resultnum, a.appname as appname, a.gtnum as gtnum, a.ival as ival, a.keyword as keyword,a.uid as uid, a.cid as cid from alarmhistory as ah, alarm as a where ah.jobid = a.id`
	_ = db.Select(&historys, sql)
	for _, history := range historys {
		hsql := `update alarmhistory set uid=:uid,cid=:cid,apname=:appname,keyword=:keyword,ival=:ival,gtnum=:gtnum where id=:id`
		_, _ = db.NamedExec(hsql, history)
	}
	log.Debug("update end...")
}

func InitDB() (*sqlx.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		config.GetConfig().Mc.UserName,
		config.GetConfig().Mc.PassWord,
		config.GetConfig().Mc.Host,
		config.GetConfig().Mc.Port,
		config.GetConfig().Mc.DataBase)
	db, err := sqlx.Open("mysql", uri)
	if err != nil {
		log.Errorf("cat not connection mysql error: %v, uri:%s", err, uri)
		return db, err
		//panic(-1)
	}
	err = db.Ping()
	if err != nil {
		log.Error("can not ping mysql error: ", err)
		return db, err
	}
	err, maxIdleConns := config.GetStringMapInt("mysql", "maxIdleConns")
	if err == nil {
		db.SetMaxIdleConns(maxIdleConns)
	}
	err, maxOpenConns := config.GetStringMapInt("mysql", "maxOpenConns")
	if err == nil {
		db.SetMaxOpenConns(maxOpenConns)
	}
	return db, err
}
