# quiver 配置中心

一个高性能、高并发的分布式配置中心服务，基于 Go Fiber 框架构建，支持 RESTful API 接口和实时配置推送。

## 🚀 特性

- **高性能**: 基于 Fiber 框架，支持高并发请求处理
- **分层架构**: App -> Cluster -> Namespace -> Item 四层资源管理
- **实时推送**: 支持Http长轮询和websocket 两种方式的配置变更通知
- **增量更新**: 支持配置的增量拉取，减少网络传输
- **版本管理**: 完整的配置发布历史和版本管理
- **API 文档**: 集成 Markdown 文档，支持在线查阅

## 📋 系统要求

- Go 1.21+
- MySQL 8.0+

## 🛠️ 安装部署

### 1. 克隆项目
```bash
git clone https://github.com/stevenrao/quiver.git
cd quiver
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库和连接信息
```

### 4. 创建数据库
```bash
# 连接 MySQL 并创建数据库
mysql -u root -p
CREATE DATABASE quiver CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 启动服务
```bash
go run main.go
```

服务启动后，可以通过以下地址访问：
- API 服务: http://localhost:8080
- API 文档: http://localhost:8080/docs/api.html
- 健康检查: http://localhost:8080/health

## 📖 API 文档

### 基础信息
- **Base URL**: `http://your-server:8080/api/v1`
- **响应格式**: JSON
- **字符编码**: UTF-8

### 统一响应格式
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 核心接口

#### 1. 应用管理
```http
# 创建应用
POST   /api/v1/envs/{env}/apps

# 获取应用列表                    
GET    /api/v1/envs/{env}/apps   

# 获取应用详情                 
GET    /api/v1/envs/{env}/apps/{app_name}  

# 更新应用       
PUT    /api/v1/envs/{env}/apps/{app_name}   

# 删除应用      
DELETE /api/v1/envs/{env}/apps/{app_name}         
```

#### 2. 集群管理
```http
# 创建集群
POST   /api/v1/envs/{env}/apps/{app_name}/clusters     

# 获取集群列表               
GET    /api/v1/envs/{env}/apps/{app_name}/clusters      

# 获取集群详情              
GET    /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}  

# 删除集群   
DELETE /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}     
```

#### 3. 命名空间管理
```http
# 创建命名空间
POST   /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces   

# 获取命名空间列表                       
GET    /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces     

# 获取命名空间详情                     
GET    /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}    

# 删除命名空间      
DELETE /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}          
```

#### 4. Item管理
```http
# 设置一个 Key-Value 对
POST   /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}/items   

# 获取一个 Key-Value                      
GET    /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces{namespace_name}/items/{key}     

# 删除一个 Key-Value                    
DELETE /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}/items/{key}    
       
```

#### 5. 配置管理
```http
# 获取命名空间完整配置（支持分页和增量更新）
GET    /api/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/configs

# 配置变更通知（长轮询）
GET    /api/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/notifications

# 发布命名空间配置
POST   /api/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/releases

```

## 🏗️ 架构设计

### 资源层级
```
Env (环境)
└── App (应用)
    └── Cluster (集群)
        └── Namespace (命名空间)
            └── Item (配置项)
```

### 数据库表结构
- `user`: 用户表
- `app`: 应用表
- `cluster`: 集群表
- `namespace`: 命名空间表
- `item`: 配置项表
- `item_release`: 发布配置项快照表
- `namespace_release`: 命名空间发布记录表

### 关键特性

#### 1. 高性能缓存
- 缓存最新发布的命名空间配置
- 单个配置项缓存，减少数据库查询
- 支持缓存失效和主动清理

#### 2. 增量更新
- 客户端提供 `releaseKey` 参数
- 服务端计算配置差异 (adds/updates/deletes)
- 减少网络传输，提升性能

#### 3. 实时通知
- 长轮询机制，减少客户端轮询频率
- 支持超时控制，避免连接占用

#### 4. 版本管理
- 每次发布生成唯一的 `releaseKey`
- 保留完整的发布历史记录
- 支持配置回滚（预留接口）

## 🔧 开发指南

### 项目结构
```
quiver/
├── main.go              # 入口文件
├── docs/                # 文档
├── config/              # 配置管理
├── database/            # 数据库连接
├── models/              # 数据模型
├── services/            # 业务逻辑层
├── handler/             # api接入处理层
├── middleware/          # 中间件
├── utils/               # 工具函数
├── logger/              # 日志
├── script/              # 脚本
├── sdk/                 # 客户端sdk
└── routes/              # 路由定义
```

### 开发规范
1. 遵循 RESTful API 设计原则
2. 使用 GORM 进行 ORM 映射
3. 统一的错误处理和响应格式
4. 输入参数验证和 SQL 注入防护
5. 使用乐观锁防止并发更新冲突

### 性能优化
1. **数据库优化**
    - 建立合适的索引
    - 使用连接池管理数据库连接
    - 分页查询避免大结果集

2. **缓存策略**
    - 热点数据缓存
    - 缓存过期时间合理设置
    - 缓存穿透和雪崩防护

3. **并发控制**
    - 使用乐观锁版本控制
    - 限流中间件防止服务过载
    - 长连接超时控制

## 📊 性能指标

- **QPS**: 支持万级别并发请求
- **延迟**: 平均响应时间 < 10ms
- **容量**: 支持万级别命名空间配置项
- **可用性**: 99.9% 服务可用性

## 🚀 部署建议

### 生产环境部署
1. 使用 Docker 容器化部署
2. 配置 MySQL 主从复制
4. 负载均衡器分发请求
5. 监控和日志收集

### 扩展性考虑
1. 支持水平扩展
2. 数据库分库分表
3. 缓存分片
4. 配置热更新

## 📄 许可证

本项目基于 GNU General Public License v3.0 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系方式

- 项目主页: https://github.com/mageg-x/quiver
- 问题反馈: https://github.com/mageg-x/quiver/issues
- 邮箱: stevenrao@me.com

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！