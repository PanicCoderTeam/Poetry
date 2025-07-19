-- 游戏房间表
CREATE TABLE `game_room` (
  -- 核心标识属性
  `id` bigint  NOT NULL AUTO_INCREMENT,
  `room_id` VARCHAR(12) NOT NULL COMMENT '短房间ID，4-7位数字',
  
  -- 基础配置属性
  `status` ENUM('waiting', 'playing', 'closed') NOT NULL DEFAULT 'waiting' COMMENT '房间状态',
  `max_players` TINYINT UNSIGNED NOT NULL DEFAULT 4 COMMENT '最大玩家数',
  `current_players` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '当前玩家数',
  `password` VARCHAR(6) DEFAULT NULL COMMENT '房间密码，6位数字',
  `min_level` SMALLINT UNSIGNED DEFAULT 0 COMMENT '最低准入等级',
  `cost_item_id` INT DEFAULT NULL COMMENT '消耗道具ID',
  `cost_amount` INT DEFAULT NULL COMMENT '消耗道具数量',
  
  -- 功能扩展属性
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `game_mode` VARCHAR(32) NOT NULL DEFAULT 'classic' COMMENT '游戏模式',
  `slogan` VARCHAR(32) DEFAULT NULL COMMENT '房间标语',
  `owner_id` INT NOT NULL COMMENT '房主ID',
  `player_list` text DEFAULT NULL COMMENT '玩家列表，存储玩家ID、准备状态等信息',
  
  -- 索引
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_room_id` (`room_id`),
  KEY `idx_status` (`status`),
  KEY `idx_owner` (`owner_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='游戏房间表';

-- 房间ID生成函数(示例)
DELIMITER //
CREATE FUNCTION generate_room_id() RETURNS VARCHAR(12)
BEGIN
  DECLARE new_id VARCHAR(12);
  DECLARE exists_flag INT;
  
  REPEAT
    -- 生成6位随机数字
    SET new_id = LPAD(FLOOR(RAND() * 999999), 6, '0');
    
    -- 检查是否已存在
    SELECT COUNT(*) INTO exists_flag FROM `game_room` WHERE `room_id` = new_id;
  UNTIL exists_flag = 0 END REPEAT;
  
  RETURN new_id;
END //
DELIMITER ;