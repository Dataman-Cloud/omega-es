use alarm;

alter table alarm modify createtime timestamp not null default CURRENT_TIMESTAMP;
