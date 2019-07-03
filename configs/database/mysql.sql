DROP TABLE IF EXISTS `kubernetes_cluster`;
CREATE TABLE IF NOT EXISTS `kubernetes_cluster`  (
      `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
      `name` varchar(16) NOT NULL DEFAULT '' COMMENT 'name',
      `region` varchar(16) NOT NULL DEFAULT '' COMMENT 'region',
      `namespace` varchar(16) NOT NULL DEFAULT '' COMMENT 'k8s namespace',
      `master_ip` varchar(16) NOT NULL DEFAULT '' COMMENT 'sub k8s cluster master',
      `master_port` varchar(255) NOT NULL DEFAULT '' COMMENT '',
      `tls_cert` varchar(255) NOT NULL DEFAULT '' COMMENT 'sub k8s tls',
      `tls_key` varchar(255) NOT NULL DEFAULT '' COMMENT '',
      `user` varchar(64) NOT NULL DEFAULT '' COMMENT 'user',
      `task_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'task id',
      `create_time` datetime NOT NULL DEFAULT '1970-01-01' COMMENT '创建时间戳',
      `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
      PRIMARY KEY (`id`),
      UNIQUE KEY `uniq_order_id` (`region`, `name`,`namespace`),
      INDEX idx_node(`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='kubernetes cluster';


DROP TABLE IF EXISTS `operation_log`;
CREATE TABLE `operation_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `owner_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '操作对象id',
  `task_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '相关的taskid',
  `message` text COMMENT '备注信息',
  `user` varchar(32) NOT NULL DEFAULT '' COMMENT '操作用户',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:01' COMMENT '触发时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=244277 DEFAULT CHARSET=utf8 COMMENT='操作记录表';
