use alarm;

alter table app_event modify createtime timestamp not null default CURRENT_TIMESTAMP;
