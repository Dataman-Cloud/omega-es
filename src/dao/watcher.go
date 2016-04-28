package dao

import (
	//"encoding/base32"
	"github.com/Dataman-Cloud/omega-es/src/model"
	"github.com/Dataman-Cloud/omega-es/src/util"
	log "github.com/cihub/seelog"
)

func ExistWatcherByUidAndUtypeAndWname(uid int64, utype, wname string) (int, error) {
	db := util.DB()
	count := 0
	sql := `select count(*) from watcher where uid = ? and utype = ? and wname = ?`
	err := db.Get(&count, sql, uid, utype, wname)
	return count, err
}

func GetWatcherByUidAndUtypeAndWname(uid int64, utype, wname string) (model.Watcher, error) {
	db := util.DB()
	watcher := model.Watcher{}
	sql := `select * from watcher where uid = ? and utype = ? and wname = ?`
	err := db.Get(&watcher, sql, uid, utype, wname)
	return watcher, err
}

func GetWatchersByUser(uid, pagecount, pagenum int64, utype string) ([]model.Watcher, error) {
	db := util.DB()
	watchers := []model.Watcher{}
	err := db.Select(&watchers, `select * from watcher where uid = ? and utype = ? limit ?, ?`, uid, utype, (pagenum-1)*pagecount, pagecount)
	return watchers, err
}

func InsertWatcher(watcher *model.Watcher) (int64, error) {
	//watcher.Cwname = base32.StdEncoding.EncodeToString([]byte(watcher.Wname))
	db := util.DB()
	sql := `insert into watcher(uid, utype, wname, wbody, cwname, notify) values (:uid, :utype, :wname, :wbody, :cwname, :notify)`
	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		log.Error("insert into watcher error: ", err)
		return 0, err
	}
	defer func() {
		if stmt != nil {
			err = stmt.Close()
			if err != nil {
				log.Error("insert into watcher close stmt error: ", err)
			}
		}
	}()
	result, err := stmt.Exec(watcher)
	if err != nil {
		log.Error("insert into watcher error: ", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func UpdateWatcher(watcher *model.Watcher) error {
	db := util.DB()
	sql := `update watcher set wbody = :wbody where uid = :uid and utype = :utype and wname = :wname`
	_, err := db.NamedExec(sql, watcher)
	return err
}

func DeleteWatcher(watcher *model.Watcher) error {
	db := util.DB()
	sql := `delete from watcher where uid = :uid and utype = :utype and wname = :wname`
	_, err := db.NamedExec(sql, watcher)
	return err
}
