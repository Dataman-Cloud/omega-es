use alarm;

alter table alarmhistory add uid bigint(20) NOT NULL;
alter table alarmhistory add cid bigint(20) NOT NULL;
alter table alarmhistory add appname varchar(64) NOT NULL;
alter table alarmhistory add keyword varchar(120) NOT NULL;
alter table alarmhistory add ival int(5) NOT NULL;
alter table alarmhistory add gtnum int(5) NOT NULL;
