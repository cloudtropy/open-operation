
CREATE TABLE IF NOT EXISTS `ops_host` (
  `host_id`       VARCHAR(32) NOT NULL COMMENT '主机唯一id',
  `ip`            VARCHAR(17) NOT NULL COMMENT 'tcp ip',
  `status`        TINYINT(1) DEFAULT 0 NOT NULL COMMENT '主机状态，可用否 0上架 1下架',
  `hostname`      VARCHAR(255) COMMENT 'hostname',
  `comment`       VARCHAR(1024) DEFAULT '',
  `location`      VARCHAR(1024) DEFAULT '',
  `cpu_count`     SMALLINT UNSIGNED NOT NULL,
  `mem_capacity`  BIGINT UNSIGNED NOT NULL,
  `disk_capacity` BIGINT UNSIGNED NOT NULL,
  `os`            VARCHAR(255) DEFAULT '' COMMENT 'Operating System',
  `create_time`   DATETIME NOT NULL DEFAULT now(),
  `update_time`   DATETIME NOT NULL DEFAULT now(),
  PRIMARY KEY (`host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='主机基础信息表';

CREATE TABLE IF NOT EXISTS `ops_host_server` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `job_id`        INT UNSIGNED NOT NULL,
  `host_id`       VARCHAR(32) NOT NULL COMMENT '主机唯一id，由agent首次启动时生成',
  `version`       VARCHAR(255) DEFAULT '',
  `weight`        INT NOT NULL DEFAULT 0,
  `create_time`   DATETIME NOT NULL DEFAULT now(),
  `update_time`   DATETIME NOT NULL DEFAULT now(),
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_job_id_host_id` (`job_id`, `host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='当前已分配的服务';

CREATE TABLE IF NOT EXISTS `ops_job` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`          VARCHAR(255) NOT NULL,
  `service_id`    INT NOT NULL COMMENT '-1代表未分配',
  `job_type`      VARCHAR(64) NOT NULL,
  `comment`       VARCHAR(1024) DEFAULT '',
  `creator`       VARCHAR(64) DEFAULT '',
  `create_time`   DATETIME NOT NULL DEFAULT now(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='';

CREATE TABLE IF NOT EXISTS `ops_job_template` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `job_id`        INT UNSIGNED NOT NULL,
  `template_id`   INT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_template_job` (`template_id`, `job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='';

CREATE TABLE IF NOT EXISTS `ops_template` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`          VARCHAR(64) NOT NULL,
  `description`   VARCHAR(1024) DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_template` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='监控模板';

CREATE TABLE IF NOT EXISTS `ops_item` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`          VARCHAR(64) NOT NULL,
  `interval`      INT UNSIGNED NOT NULL COMMENT '汇报间隔,单位秒',
  `history`       INT UNSIGNED NOT NULL COMMENT '历史数据存储周期,单位天',
  `description`   VARCHAR(1024) DEFAULT '',
  `data_type`     VARCHAR(64) NOT NULL,
  `unit`          VARCHAR(64) NOT NULL,
  `creator`       VARCHAR(64) NOT NULL DEFAULT 'born/system/(user)',
  `dst`           VARCHAR(64) NOT NULL COMMENT 'GAUGE/COUNTER',
  `create_time`   DATETIME NOT NULL DEFAULT now(),
  `update_time`   DATETIME NOT NULL DEFAULT now(),
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='监控项目表';

CREATE TABLE IF NOT EXISTS `ops_template_item` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `template_id`   INT UNSIGNED NOT NULL,
  `item_id`       INT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_template_item` (`template_id`, `item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='';

CREATE TABLE IF NOT EXISTS `ops_trigger` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`          VARCHAR(64) NOT NULL,
  `severity`      VARCHAR(16) NOT NULL COMMENT '',
  `notice_message`  VARCHAR(10240) DEFAULT '',
  `notice_person` VARCHAR(512) NOT NULL,
  `notice_by`     VARCHAR(64) NOT NULL,
  `template_id`   INT NOT NULL DEFAULT -1 COMMENT '-1代表未分配',
  `item_id`       INT UNSIGNED NOT NULL,
  `rule_time`     VARCHAR(16) NOT NULL,
  `rule_type`     VARCHAR(16) NOT NULL,
  `rule_operator` VARCHAR(16) NOT NULL,
  `rule_value`    VARCHAR(32) NOT NULL,
  `enabled`       TINYINT(1) NOT NULL DEFAULT 0 COMMENT '0 enabled 1 not enabled',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_template_name` (`template_id`, `name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='监控告警信息表';

CREATE TABLE IF NOT EXISTS `ops_trigger_msg` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `trigger_id`    INT UNSIGNED NOT NULL,
  `host_id`       VARCHAR(32) NOT NULL,
  `status`        TINYINT(1) NOT NULL DEFAULT 1 COMMENT '0 ok 1 problem',
  `ack`           TINYINT(1) NOT NULL DEFAULT 1 COMMENT '0 yes 1 no',
  `ack_msg`       VARCHAR(512) DEFAULT '',
  `ack_user`      VARCHAR(64) DEFAULT '',
  `ack_time`      DATETIME DEFAULT now(),
  `create_time`   DATETIME NOT NULL DEFAULT now(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='';

CREATE TABLE IF NOT EXISTS `ops_graph` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`          VARCHAR(64) NOT NULL,
  `description`   VARCHAR(1024) DEFAULT '',
  `template_id`   INT NOT NULL DEFAULT -1 COMMENT '-1代表未分配',
  `graph_type`    VARCHAR(64) NOT NULL,
  `host_id`       VARCHAR(32) NOT NULL COMMENT '主机表id',
  `screen_id`     INT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_graph` (`name`, `host_id`, `screen_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='';


CREATE TABLE IF NOT EXISTS `ops_graph_item` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `graph_id`      INT UNSIGNED NOT NULL,
  `item_id`       INT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_graph_item` (`graph_id`, `item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='';

CREATE TABLE IF NOT EXISTS `ops_screen` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`          VARCHAR(64) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_screen` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='';



CREATE TABLE IF NOT EXISTS `ops_user` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user`          VARCHAR(64) NOT NULL,
  `alias`         VARCHAR(12) NOT NULL DEFAULT '',
  `passwd`        CHAR(32) NOT NULL COMMENT 'user password',
  `salt`          CHAR(24) NOT NULL COMMENT 'hash salt of password',
  `email`         VARCHAR(255) NOT NULL DEFAULT '',
  `phone`         VARCHAR(11) NOT NULL DEFAULT '',
  `wechat`        VARCHAR(64) NOT NULL DEFAULT '',
  `sex`           VARCHAR(6) NOT NULL COMMENT 'male or female',
  `create_time`   DATETIME NOT NULL DEFAULT now(),
  `role`          TINYINT NOT NULL DEFAULT 0 COMMENT '-1:blocked 0:normal 1:root',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `user` (`user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户信息';

-- root
INSERT INTO `ops_user` (`id`, `user`, `alias`, `passwd`, `salt`, `sex`, `role`) VALUES(1, 'root', '管理员', '1374BCE70F8095EACC3DCEE34C894DEE', '6f180a2d979246b7bb580e00', 'male', 1);

CREATE TABLE IF NOT EXISTS `ops_team_user` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tid`           INT UNSIGNED NOT NULL,
  `uid`           INT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_team_user` (`uid`, `tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户与组';

CREATE TABLE IF NOT EXISTS `ops_action_trail` (
  `id`            VARCHAR(32) NOT NULL,
  `user`          VARCHAR(64) NOT NULL,
  `module`        VARCHAR(64) NOT NULL,
  `action`        VARCHAR(64) NOT NULL,
  `result`        VARCHAR(512) NOT NULL DEFAULT '',
  `detail`        VARCHAR(1024) NOT NULL DEFAULT '',
  `create_time`   DATETIME NOT NULL DEFAULT now(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='操作审计';


INSERT INTO `ops_item` VALUES ('1', 'cpu.load.1min', 60, 365, 'CPU load 1min. Returns float', 'float', '', 'system', 'GAUGE', now(), now());
INSERT INTO `ops_item` VALUES ('2', 'cpu.load.5min', 60, 365, 'CPU load 5min. Returns float', 'float', '', 'system', 'GAUGE', now(), now());
INSERT INTO `ops_item` VALUES ('3', 'cpu.load.15min', 60, 365, 'CPU load 15min. Returns float', 'float', '', 'system', 'GAUGE', now(), now());
INSERT INTO `ops_item` VALUES ('4', 'cpu.load.1min.percent', 60, 365, 'CPU load 1min percent. Returns float percent[0, 100.0]', 'float', '%', 'system', 'GAUGE', now(), now());
INSERT INTO `ops_item` VALUES ('5', 'mem.free', 60, 365, 'Memory free. Returns uint', 'uint', 'kB', 'system', 'GAUGE', now(), now());
INSERT INTO `ops_item` VALUES ('6', 'mem.used', 60, 365, 'Memory used. Returns uint', 'uint', 'kB', 'system', 'GAUGE', now(), now());
INSERT INTO `ops_item` VALUES ('7', 'mem.used.percent', 60, 365, 'Memory used percent. Returns float percent[0, 100.0]', 'float', '%', 'system', 'GAUGE', now(), now());
INSERT INTO `ops_item` VALUES ('8', 'net.if.in', 60, 365, '', 'float', 'B/s', 'system', 'COUNTER', now(), now());
INSERT INTO `ops_item` VALUES ('9', 'net.if.out', 60, 365, '', 'float', 'B/s', 'system', 'COUNTER', now(), now());
INSERT INTO `ops_item` VALUES ('10', 'net.if.total', 60, 365, '', 'float', 'B/s', 'system', 'COUNTER', now(), now());
INSERT INTO `ops_item` VALUES ('11', 'df.free', 60, 365, '', 'uint', 'kB', 'system', 'GAUGE', now(), now());
INSERT INTO `ops_item` VALUES ('12', 'df.used', 60, 365, '', 'uint', 'kB', 'system', 'GAUGE', now(), now());
INSERT INTO `ops_item` VALUES ('13', 'df.used.percent', 60, 365, '', 'float', '%', 'system', 'GAUGE', now(), now());
