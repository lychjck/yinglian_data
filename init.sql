-- 创建数据库
CREATE DATABASE IF NOT EXISTS yinglian_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE yinglian_db;

-- 创建古籍内容表
CREATE TABLE ancient_books (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    book_name VARCHAR(100)  COMMENT '书名',
    volume VARCHAR(100)  COMMENT '卷数',
    title VARCHAR(100) COMMENT '标题',
    content TEXT NOT NULL COMMENT '内容',
    ref INT COMMENT '关联楹联',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='古籍内容表';
