package model

import (
	"time"
)

type Watcher struct {
	Id     int64  `json:"id"`
	Uid    int64  `json:"uid"`
	Utype  string `json:"utype"`
	Wname  string `json:"wname"`
	Cwname string `json:"cwname"`
	Wbody  string `json:"wbody"`
	Notify bool   `json:"notify"`
}

type CreateWatcher struct {
	Id      string `json:"_id"`
	Version int    `json:"_version"`
	Created bool   `json:"created"`
}

type LogAlarm struct {
	Id         int64     `json:"id"`
	Uid        int64     `json:"uid"`
	Cid        int64     `json:"cid"`
	AppName    string    `json:"appname"`
	AppAlias   string    `json:"appalias"`
	Ival       int8      `json:"ival"`
	GtNum      int64     `json:"gtnum"`
	AlarmName  string    `json:"alarmname"`
	UserType   string    `json:"usertype"`
	KeyWord    string    `json:"keyword"`
	Emails     string    `json:"emails"`
	AliasName  string    `json:"aliasname"`
	CreateTime time.Time `json:"createtime"`
}

type AlarmHistory struct {
	Id        int64     `json:"id"`
	JobId     int64     `json:"jobid"`
	IsAlarm   bool      `json:"isalarm"`
	ExecTime  time.Time `json:"exectime"`
	ResultNum int64     `json:"resultnum"`
	Ival      int8      `json:"ival"`
	GtNum     int64     `json:"gtnum"`
	AppName   string    `json:"appname"`
	KeyWord   string    `json:"keyword"`
	Uid       int64     `json:"uid"`
	Cid       int64     `json:"cid"`
}
