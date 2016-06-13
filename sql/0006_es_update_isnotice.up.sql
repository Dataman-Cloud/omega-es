use alarm;

alter table alarm modify isnotice tinyint(1) not null default 1;
