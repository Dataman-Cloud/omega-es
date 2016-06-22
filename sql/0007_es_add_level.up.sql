use alarm;
alter table alarm add level varchar(8) not null default "";
alter table alarmhistory add level varchar(8) not null default "";
