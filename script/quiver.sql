-- MySQL 脚本：KVConfig 配置中心数据库表结构
-- 文件编码：UTF-8
-- 适用于 MySQL 8+

-- 如果数据库已存在，可选择删除后重建（生产环境慎用）
-- DROP DATABASE IF EXISTS quiver_pro;
-- CREATE DATABASE quiver_pro CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE quiver;

-- 用户表
-- 用户表
CREATE TABLE IF NOT EXISTS user (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_name   VARCHAR(128) NOT NULL,
    password    VARCHAR(256) NOT NULL,
    email       VARCHAR(128),
    phone       VARCHAR(32),
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_user_name (user_name),
    UNIQUE KEY uk_email (email),
    UNIQUE KEY uk_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 应用表
CREATE TABLE IF NOT EXISTS app (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    app_name    VARCHAR(128) NOT NULL,
    description VARCHAR(1024),
    ver         BIGINT NOT NULL DEFAULT 1,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_app_name (app_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 集群表
CREATE TABLE IF NOT EXISTS cluster (
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    app_id       BIGINT NOT NULL,
    app_name     VARCHAR(128) NOT NULL,
    cluster_name VARCHAR(128) NOT NULL,
    description  VARCHAR(1024),
    ver          BIGINT NOT NULL DEFAULT 1,
    create_time  DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_app_cluster_name (app_name, cluster_name),
    FOREIGN KEY (app_id) REFERENCES app (id) ON DELETE CASCADE,

    -- 优化：合并查询路径，支持 app_id + cluster_name 查询
    KEY idx_app_cluster (app_id, cluster_name),
    KEY idx_app_name (app_name),
    KEY idx_cluster_name (cluster_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 命名空间表
CREATE TABLE IF NOT EXISTS namespace (
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    namespace_name VARCHAR(128) NOT NULL,
    description    VARCHAR(1024),
    app_id         BIGINT NOT NULL,
    app_name       VARCHAR(128) NOT NULL,
    cluster_id     BIGINT NOT NULL,
    cluster_name   VARCHAR(128) NOT NULL,
    ver            BIGINT NOT NULL DEFAULT 1,
    create_time    DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_app_cluster_namespace (app_name, cluster_name, namespace_name),
    FOREIGN KEY (app_id) REFERENCES app (id) ON DELETE CASCADE,
    FOREIGN KEY (cluster_id) REFERENCES cluster (id) ON DELETE CASCADE,

    -- 优化：高频查询路径
    KEY idx_app_cluster_namespace (app_id, cluster_id, namespace_name),
    KEY idx_cluster_namespace (cluster_id, namespace_name),
    KEY idx_namespace_name (namespace_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 配置项表 (item)
CREATE TABLE IF NOT EXISTS item (
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    app_id         BIGINT NOT NULL,
    cluster_id     BIGINT NOT NULL,
    namespace_id   BIGINT NOT NULL,
    k              VARCHAR(255) NOT NULL,
    v              TEXT,
    kv_id          BIGINT UNSIGNED NOT NULL,

    create_time    DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_deleted     TINYINT DEFAULT 0,
    is_released    TINYINT DEFAULT 0,

    UNIQUE KEY uk_namespace_k (namespace_id, k),

    -- 优化：合并高频查询
    KEY idx_namespace_k (namespace_id, k),           -- 查询 k
    KEY idx_namespace_released (namespace_id, is_released, is_deleted), -- 发布查询
    KEY idx_namespace_deleted (namespace_id, is_deleted), -- 批量操作
    KEY idx_kv_id (kv_id),                           -- 快速查 kv_id
    KEY idx_k (k)                                    -- 全局查 key（可选）
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 命名空间发布记录表
-- 命名空间发布记录表
CREATE TABLE IF NOT EXISTS namespace_release (
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    app_id         BIGINT NOT NULL,
    app_name       VARCHAR(128) NOT NULL,
    cluster_id     BIGINT NOT NULL,
    cluster_name   VARCHAR(128) NOT NULL,
    namespace_id   BIGINT NOT NULL,
    namespace_name VARCHAR(128) NOT NULL,
    release_id     VARCHAR(64) NOT NULL,
    release_name   VARCHAR(128) NOT NULL,
    release_time   DATETIME DEFAULT CURRENT_TIMESTAMP,
    operator       VARCHAR(64),
    comment        VARCHAR(1024),
    config         BLOB,
    create_time    DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_release_id (release_id),
    FOREIGN KEY (app_id) REFERENCES app (id) ON DELETE CASCADE,
    FOREIGN KEY (cluster_id) REFERENCES cluster (id) ON DELETE CASCADE,

    -- 优化：高频查询
    KEY idx_namespace_time (namespace_id, release_time DESC), -- 最近发布
    KEY idx_cluster_namespace (cluster_id, namespace_name),   -- 按集群+命名空间查
    KEY idx_release_time (release_time DESC)                  -- 全局发布时间
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 发布配置项表
CREATE TABLE IF NOT EXISTS item_release (
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    app_id         BIGINT NOT NULL,
    cluster_id     BIGINT NOT NULL,
    namespace_id   BIGINT NOT NULL,
    k              VARCHAR(255) NOT NULL,
    v              TEXT,
    kv_id          BIGINT UNSIGNED NOT NULL,
    is_deleted     TINYINT DEFAULT 0,

    create_time    DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_namespace_kv_id (namespace_id, kv_id),

    -- 优化：与 item 表保持一致
    KEY idx_namespace_deleted (namespace_id, is_deleted),
    KEY idx_kv_id (kv_id),
    KEY idx_namespace_k (namespace_id, k)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 可选：添加注释说明
ALTER TABLE user COMMENT '用户表';
ALTER TABLE app COMMENT '应用表';
ALTER TABLE cluster COMMENT '集群表';
ALTER TABLE namespace COMMENT '命名空间表';
ALTER TABLE item COMMENT '配置项表';
ALTER TABLE item_release COMMENT '配置项发布快照表';
ALTER TABLE namespace_release COMMENT '命名空间发布记录表';
