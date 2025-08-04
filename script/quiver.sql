-- MySQL 脚本：KVConfig 配置中心数据库表结构
-- 文件编码：UTF-8
-- 适用于 MySQL 8+

-- 如果数据库已存在，可选择删除后重建（生产环境慎用）
DROP DATABASE IF EXISTS quiver_pro;
CREATE DATABASE quiver_pro CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE quiver_pro;

-- 用户表
CREATE TABLE user
(
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_name   VARCHAR(128) NOT NULL,
    password    VARCHAR(256) NOT NULL, -- 建议存储加密后的密码（如 bcrypt）
    email       VARCHAR(128),
    phone       VARCHAR(32),
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_user_name (user_name),
    UNIQUE KEY uk_email (email),
    UNIQUE KEY uk_user_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 应用表
CREATE TABLE app
(
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    app_name    VARCHAR(128) NOT NULL,
    description VARCHAR(1024),
    ver         BIGINT      NOT NULL DEFAULT 1, -- CAS 乐观锁
    create_time DATETIME    DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- 唯一约束
    UNIQUE KEY uk_app_name (app_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 集群表
CREATE TABLE cluster
(
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    app_id       BIGINT  NOT NULL,
    app_name     VARCHAR(128) NOT NULL,
    cluster_name VARCHAR(128) NOT NULL,
    description  VARCHAR(1024),
    ver          BIGINT     NOT NULL DEFAULT 1,
    create_time  DATETIME   DEFAULT CURRENT_TIMESTAMP,
    update_time  DATETIME   DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- 唯一性：一个 app 下 cluster_name 唯一
    UNIQUE KEY uk_app_cluster_name (app_name, cluster_name),
    -- 外键（可选，MySQL 中可用）
    FOREIGN KEY (app_id) REFERENCES app (id) ON DELETE CASCADE,
    -- 查询加速
    KEY          idx_app_id (app_id),
    KEY          idx_app_name (app_name),
    KEY          idx_cluster_name (cluster_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 命名空间表
CREATE TABLE namespace
(
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    namespace_name VARCHAR(128) NOT NULL,
    description    VARCHAR(1024),
    app_id         BIGINT NOT NULL,
    app_name       VARCHAR(128) NOT NULL,
    cluster_id     BIGINT NOT NULL,
    cluster_name   VARCHAR(128) NOT NULL,
    ver            BIGINT       NOT NULL DEFAULT 1,
    create_time    DATETIME     DEFAULT CURRENT_TIMESTAMP,
    update_time    DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- 唯一性：同一 cluster 下 namespace_name 唯一
    UNIQUE KEY uk_app_cluster_namespace (app_name, cluster_name, namespace_name),
    -- 外键
    FOREIGN KEY (app_id) REFERENCES app (id) ON DELETE CASCADE,
    FOREIGN KEY (cluster_id) REFERENCES cluster (id) ON DELETE CASCADE,

    -- 查询加速
    KEY            idx_app_id (app_id),
    KEY            idx_cluster_id (cluster_id),
    KEY            idx_app_name (app_name),
    KEY            idx_cluster_name (cluster_name),
    KEY            idx_namespace_name (namespace_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 配置项表
CREATE TABLE item
(
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    app_id         BIGINT NOT NULL,
    cluster_id     BIGINT NOT NULL,
    namespace_id   BIGINT NOT NULL,
    k              VARCHAR(255) NOT NULL,
    v              TEXT,
    namespace_name VARCHAR(128) NOT NULL,
    ver            BIGINT       NOT NULL DEFAULT 1,
    create_time    DATETIME     DEFAULT CURRENT_TIMESTAMP,
    update_time    DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted        TINYINT      DEFAULT NULL,

    -- 唯一性
    UNIQUE KEY uk_namespace_key (namespace_id, k),

    -- 查询加速
    KEY            idx_app_id (app_id),
    KEY            idx_cluster_id (cluster_id),
    KEY            idx_namespace_id (namespace_id),
    KEY            idx_namespace_name (namespace_name),
    KEY            idx_k (k)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 命名空间发布记录表
CREATE TABLE namespace_release
(
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    app_id         BIGINT NOT NULL,
    app_name       VARCHAR(128) NOT NULL,
    cluster_id     BIGINT NOT NULL,
    cluster_name   VARCHAR(128) NOT NULL,
    namespace_id   BIGINT NOT NULL,
    namespace_name VARCHAR(128) NOT NULL,
    release_id     VARCHAR(64)  NOT NULL,
    release_time   DATETIME DEFAULT CURRENT_TIMESTAMP,
    operator       VARCHAR(64),
    comment        VARCHAR(1024),
    config         BLOB,    -- 本次 release 包含的 item 的 (kv_id)，MessagePack 编码
    create_time    DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- 唯一性
    UNIQUE KEY uk_release_id (release_id),

    -- 外键
    FOREIGN KEY (app_id) REFERENCES app (id) ON DELETE CASCADE,
    FOREIGN KEY (cluster_id) REFERENCES cluster (id) ON DELETE CASCADE,

    -- 查询加速（最重要）
    KEY            idx_namespace_id (namespace_id),
    KEY            idx_namespace_name (namespace_name),
    KEY            idx_namespace_time (namespace_name, release_time DESC), -- 查最近发布
    KEY            idx_release_time (release_time DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 发布配置项表（item 的快照）
CREATE TABLE item_release
(
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    app_id         BIGINT NOT NULL,
    cluster_id     BIGINT NOT NULL,
    namespace_id   BIGINT NOT NULL,
    release_id     VARCHAR(64)  NOT NULL,   -- namespace_release 的 release_id
    k              VARCHAR(255) NOT NULL,
    v              TEXT,
    kv_id          BIGINT   NOT NULL,   -- MurmurHash64(k,v) 得到
    deleted        TINYINT      DEFAULT NULL,

    -- 唯一性
    UNIQUE KEY uk_namespace_kv_id (namespace_id, kv_id), -- 核心去重

    -- 索引
    KEY            idx_app_id (app_id),
    KEY            idx_cluster_id (cluster_id),
    KEY            idx_namespace_id (namespace_id),
    KEY            idx_kv_id (kv_id),
    KEY            idx_release_id (release_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 可选：添加注释说明
ALTER TABLE user COMMENT '用户表';
ALTER TABLE app COMMENT '应用表';
ALTER TABLE cluster COMMENT '集群表';
ALTER TABLE namespace COMMENT '命名空间表';
ALTER TABLE item COMMENT '配置项表';
ALTER TABLE item_release COMMENT '配置项发布快照表';
ALTER TABLE namespace_release COMMENT '命名空间发布记录表';
