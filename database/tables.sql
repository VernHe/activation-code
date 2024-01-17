create table activation_attempts
(
    id            int auto_increment
        primary key,
    card_value    varchar(255)                       not null comment '关联的激活码ID',
    activation_at datetime                           not null comment '激活尝试时间',
    success       tinyint(1)                         not null comment '是否成功激活 (1表示成功，0表示失败)',
    error_message varchar(255)                       null comment '激活失败时的错误信息',
    request_data  text                               null comment '激活请求的JSON数据',
    response_data text                               null comment '激活响应的JSON数据',
    created_at    datetime default CURRENT_TIMESTAMP not null comment '记录创建时间'
);

create table app
(
    id          varchar(32)                         not null
        primary key,
    name        varchar(255)                        not null,
    card_length int                                 not null,
    card_prefix varchar(255)                        not null,
    created_at  timestamp default CURRENT_TIMESTAMP not null
);

create table card
(
    id         varchar(36)   not null comment '卡密唯一标识符(UUID)'
        primary key,
    status     int           not null comment '状态: 0-未使用, 1-已使用, 2-已锁定, 3-已删除',
    user_id    varchar(36)   not null comment '用户ID，关联到用户表中的id字段',
    user_name  varchar(255)  null comment '用户名称',
    days       int default 0 not null comment '有效天数',
    expired_at datetime      null comment '过期时间',
    value      varchar(255)  not null comment '卡密值',
    used       tinyint(1)    not null comment '是否使用过',
    seid       varchar(255)  null comment 'SEID码',
    remark     varchar(255)  null comment '备注信息',
    used_at    datetime      null,
    deleted_at datetime      null comment '删除时间',
    created_at datetime      not null comment '生成时间',
    app_id     varchar(36)   null comment 'App唯一标识符(UUID)',
    minutes    int default 0 not null comment '有效分钟数',
    locked_at  datetime      null comment '锁定时间',
    time_type  varchar(50)   not null comment '激活码时间类型'
);

create table user
(
    id           varchar(32)                         not null
        primary key,
    username     varchar(255)                        not null,
    password     varchar(255)                        not null,
    status       int       default 1                 null,
    ancestry     varchar(255)                        null,
    total_cnt    int       default 0                 null,
    used_cnt     int       default 0                 null,
    noused_cnt   int       default 0                 null,
    deleted_cnt  int       default 0                 null,
    locked_cnt   int       default 0                 null,
    created_at   timestamp default CURRENT_TIMESTAMP not null,
    roles        json                                null,
    introduction varchar(255)                        null,
    avatar       varchar(255)                        null,
    max_cnt      int       default 0                 null,
    permissions  json                                null,
    apps         json                                null,
    constraint username
        unique (username)
);

create table user_config
(
    id           varchar(255)                        not null
        primary key,
    user_id      varchar(255)                        not null,
    config_key   varchar(255)                        not null,
    config_value text                                null,
    created_at   timestamp default CURRENT_TIMESTAMP not null,
    updated_at   timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    is_deleted   tinyint   default 0                 null,
    constraint user_id
        unique (user_id, config_key)
)
    charset = latin1;
