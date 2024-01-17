CREATE TABLE `activation_attempts`
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    card_id       VARCHAR(36)   NOT NULL COMMENT '关联的激活码ID',
    activation_at DATETIME      NOT NULL COMMENT '激活尝试时间',
    success       TINYINT(1)    NOT NULL COMMENT '是否成功激活 (1表示成功，0表示失败)',
    error_message VARCHAR(255)  NULL COMMENT '激活失败时的错误信息',
    request_data  JSON          NULL COMMENT '激活请求的JSON数据',
    response_data JSON          NULL COMMENT '激活响应的JSON数据',
    created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    FOREIGN KEY (card_id) REFERENCES configmanagement.card(id)
);