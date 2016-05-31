use alarm;

alter table alarm add ipport varchar(64) NOT NULL;
alter table alarm add scaling tinyint(1) NOT NULL DEFAULT false;
alter table alarm add maxs int(5) NOT NULL DEFAULT 0;
alter table alarm add mins int(5) NOT NULL DEFAULT 0;
alter table alarm add appid bigint(20) NOT NULL;

alter table alarmhistory add ipport varchar(64) NOT NULL;
alter table alarmhistory add scaling tinyint(1) NOT NULL DEFAULT false;
alter table alarmhistory add maxs int(5) NOT NULL DEFAULT 0;
alter table alarmhistory add mins int(5) NOT NULL DEFAULT 0;
