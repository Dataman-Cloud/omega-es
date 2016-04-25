use oapp;

/*CREATE TABLE if not exists watcher (
	  id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	  uid varchar(64) NOT NULL,
	  utype varchar(64) NOT NULL,
	  wname varchar(64) NOT NULL,
	  wbody text,
	  cwname varchar(65) NOT NULL,
	  wemails text,
	  notify tinyint(1) NOT NULL,
	  PRIMARY KEY (id),
	  UNIQUE KEY uid (uid,utype,wname)
);*/

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
