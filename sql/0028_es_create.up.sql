create database if not exists datamanalarm;
use datamanalarm;

CREATE TABLE if not exists alarm (
	id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
	uid bigint(20) NOT NULL,
	cid bigint(20) NOT NULL,
	appname varchar(64) NOT NULL,
	appalias varchar(64) NOT NULL,
	ival int(5) NOT NULL,
	gtnum int(5) NOT NULL,
	alarmname varchar(64) NOT NULL,
	usertype varchar(10) NOT NULL,
	keyword varchar(120) NOT NULL,
	emails varchar(120) NOT NULL,
	aliasname varchar(120) NOT NULL,
	createtime timestamp,
	UNIQUE KEY uniquekey (uid,usertype,alarmname)
);

CREATE TABLE if not exists alarmhistory (
	id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
	jobid bigint(20) NOT NULL,
	isalarm tinyint(1) NOT NULL,
	exectime timestamp,
	resultnum bigint(20) NOT NULL
);
