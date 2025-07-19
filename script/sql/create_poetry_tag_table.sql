CREATE TABLE `poetry_tag` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `poetry_id` bigint NOT NULL COMMENT '诗词ID',
  `tag` varchar(255) DEFAULT "" NOT NULL COMMENT '标签',
  `category` varchar(255) DEFAULT "" NOT NULL COMMENT '标签分类',
  `tag_id` bigint NOT NULL COMMENT '标签ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_poetry_tag` (`poetry_id`,`tag`, `category`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='飞花令诗词标签表';

create table `tag` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(255) DEFAULT "" NOT NULL COMMENT '标签名',
  `category` varchar(255) DEFAULT "" NOT NULL COMMENT '标签分类',
  `level` int DEFAULT 0 NOT NULL COMMENT '标签级别',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tag_name` (`name`, `category`)
)ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='标签表';