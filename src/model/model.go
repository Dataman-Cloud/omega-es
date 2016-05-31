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
	AppId      int64     `json:"appid"`
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
	Isnotice   bool      `json:"isnotice"`
	Ipport     string    `json:"ipport"`
	Scaling    bool      `json:"scaling"`
	Maxs       int8      `json:"maxs"`
	Mins       int8      `json:"mins"`
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
	Ipport    string    `json:"ipport"`
	Scaling   bool      `json:"scaling"`
	Maxs      int8      `json:"maxs"`
	Mins      int8      `json:"mins"`
}

type CronInfo struct {
	GtNum     int64  `json:"gtnum"`
	JobId     int64  `json:"jobid"`
	UserId    int64  `json:"userid"`
	ClusterId int64  `json:"clusterid"`
	KeyWord   string `json:"keyword"`
	AppName   string `json:"appname"`
	AppAlias  string `json:"appalias"`
	Ival      int8   `json:"interval"`
	UserType  string `json:"usertype"`
	AlarmName string `json:"alarmname"`
	Ipport    string `json:"ipport"`
}
