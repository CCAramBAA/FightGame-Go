-- FightGame 数据库初始化脚本
CREATE DATABASE IF NOT EXISTS fight_game DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE fight_game;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) DEFAULT '',
    avatar VARCHAR(255) DEFAULT '',
    win_count INT DEFAULT 0,
    lose_count INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 游戏房间表
CREATE TABLE IF NOT EXISTS game_rooms (
    id VARCHAR(36) PRIMARY KEY,
    room_name VARCHAR(50) NOT NULL,
    host_id BIGINT UNSIGNED NOT NULL,
    guest_id BIGINT UNSIGNED,
    status VARCHAR(20) DEFAULT 'waiting' COMMENT 'waiting/playing/finished',
    max_players INT DEFAULT 2,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_host_id (host_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 对战记录表
CREATE TABLE IF NOT EXISTS battle_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    room_id VARCHAR(36) NOT NULL,
    winner_id BIGINT UNSIGNED NOT NULL,
    loser_id BIGINT UNSIGNED NOT NULL,
    duration INT DEFAULT 0 COMMENT '对局时长(秒)',
    replay TEXT COMMENT '回放数据JSON',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_room_id (room_id),
    INDEX idx_winner_id (winner_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入测试管理员账号（密码为加密后的 admin123，需配合实际加密算法替换）
-- INSERT INTO users (username, password, nickname) VALUES ('admin', '$2a$10$...', '管理员');
