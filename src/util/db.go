package util

import (
	"fmt"
	"sync"

	"github.com/Dataman-Cloud/omega-es/src/config"
	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattes/migrate/driver/mysql"
	"github.com/mattes/migrate/migrate"
)

func MysqlInit() {
	DB()
	upgradeDB()
}

var db *sqlx.DB

func DB() *sqlx.DB {
	var err error
	if db != nil {
		return db
	}
	mutex := sync.Mutex{}
	mutex.Lock()
	db, err = InitDB()
	if err != nil {
		log.Errorf("init db error: %v", err)
		log.Flush()
		panic(-1)
	}
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
		panic(-1)
	}
	if !ok {
		log.Error("can't upgrade db")
		log.Flush()
		panic(-1)
	}
	log.Info("DB upgraded")
	log.Flush()
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
	}
	err = db.Ping()
	if err != nil {
		log.Error("can not ping mysql error: ", err)
		return db, err
	}
	db.SetMaxIdleConns(int(config.GetConfig().Mc.MaxIdleConns))
	db.SetMaxOpenConns(int(config.GetConfig().Mc.MaxOpenConns))
	return db, err
}
