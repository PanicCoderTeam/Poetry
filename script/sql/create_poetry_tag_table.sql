CREATE TABLE `poetry_tag` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `poetry_id` bigint NOT NULL COMMENT '诗词ID',
  `tag` varchar(255) DEFAULT "" NOT NULL COMMENT '标签 xx_xx_',
  `min_tag_id` bigint NOT NULL COMMENT '最小粒度标签ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_poetry_tag` (`poetry_id`,`tag`, `min_tag_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='飞花令诗词标签表';

create table `tag` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(255) DEFAULT "" NOT NULL COMMENT '标签名',
  `parent_tag` varchar(255) DEFAULT "" NOT NULL COMMENT '父级标签',
  `level` int DEFAULT 0 NOT NULL COMMENT '标签级别',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tag_name` (`name`, `parent_tag`)
)ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='标签表';

CREATE TABLE `author` (
  
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(255) DEFAULT '' NOT NULL COMMENT '作者名称',
  `desc` text NOT NULL COMMENT '赏析',
  `dynasty` varchar(128) DEFAULT '' NOT NULL COMMENT '朝代',
   PRIMARY KEY (`id`),
   UNIQUE KEY `idx_title_author` (`name`, `dynasty`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='作者表';

CREATE TABLE `poetry` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `title` varchar(255) DEFAULT '' NOT NULL COMMENT '标签名',
  `title_tradition` varchar(255) DEFAULT '' NOT NULL COMMENT '标签名_繁体',
  `paragraphs` text NOT NULL COMMENT '古诗正文',
  `paragraphs_tradition` text NOT NULL COMMENT '古诗正文_繁体',
  `author` varchar(255) DEFAULT '' NOT NULL COMMENT '作者',
  `author_tradition` varchar(255) DEFAULT '' NOT NULL COMMENT '作者_繁体',
  `dynasty` varchar(128) DEFAULT '' NOT NULL COMMENT '朝代',
  `notes` text NOT NULL COMMENT '注释',
  `comment` text NOT NULL COMMENT '赏析',
  `translation` text NOT NULL COMMENT '翻译',
  `pinyin` text NOT NULL COMMENT '拼音',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_title_author` (`title`, `paragraphs`, `author`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='诗词表';