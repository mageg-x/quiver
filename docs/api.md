## 📄 KVConfig 配置中心 RESTful API 文档（v1）

> **版本**: v1
> **Base URL**: `{config_server_url}/api/v1`
> **当前时间**: 2025年8月1日
> **设计原则**: 遵循 RESTful 风格，资源路径清晰，动词与资源分离，支持分页、增量更新与灰度预留。

---

### 🔧 通用说明

#### ✅ 响应格式（统一结构）
所有接口返回 JSON 格式，结构如下：

```json
{
  "code": 0,
  "message": "success",
  "data": {  }
}
```

- `code = 0` 表示成功
- `code ≠ 0` 表示错误，`message` 提供错误描述
- `data` 为返回的具体内容

---

### 📚 接口列表

#### 1. 获取 Namespace 最新配置（支持分页 & Delta 差异）
> 客户端拉取配置的核心接口，支持全量或增量更新。

- **URL**:
  `GET /api/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/configs`

- **Path 参数**:
    
    | 参数 | 必选 | 类型 | 说明 |
    | --- | --- | --- | --- |
    | env | 是 | string | 环境 |
    | appId | 是 | string | 应用ID |
    | clusterName | 是 | string | 集群名称 |
    | namespaceName | 是 | string | 命名空间名称 |

- **Query 参数**:

  | 参数 | 必选 | 类型 | 默认值 | 说明 |
  | --- | --- | --- | --- | --- |
  | releaseKey | 否 | string | - | 客户端当前版本，用于计算增量（delta） |
  | page | 否 | int | 1 | 页码（从1开始） |
  | size | 否 | int | 100 | 每页数量，最大 500 |
  | label | 否 | string | - | 预留：用于灰度标签 |
  | ip | 否 | string | - | 预留：用于灰度IP匹配 |

- **请求示例**:
    ```http
    GET /api/v1/envs/pro/apps/app123/clusters/default/namespaces/application/configs?releaseKey=rel-abc123&page=1&size=50
    ```

- **响应示例（200 OK）**:
    ```json
    {
      "code": 0,
      "message": "success",
      "data": {
        "env": "pro",
        "appId": "app123",
        "clusterName": "default",
        "namespaceName": "mysql",
        "releaseKey": "rel-def456",
        "comment": "数据库配置",
        "total": 150,
        "page": 1,
        "size": 50,
        "items": {
          "db.url": "jdbc:mysql://pro",
          "log.level": "INFO",
          "debug.username": "stevenrao",
          "debug.password": "stevenrao"
        },
        "changes": {
          "updates": ["db.url", "log.level"],
          "adds": ["db.username", "db.password"],
          "deletes": ["debug.mode"]
        }
      }
    }
    ```

---

#### 2. 配置变更通知（长轮询）
> 客户端长连接等待配置变化，用于实时感知更新。

- **URL**:
  `GET /api/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/notifications`

- **Path 参数**:
    
    | 参数 | 必选 | 类型 | 说明 |
    | --- | --- | --- | --- |
    | env | 是 | string | 环境 |
    | appId | 是 | string | 应用ID |
    | clusterName | 是 | string | 集群名称 |
    | namespaceName | 是 | string | 命名空间名称 |

- **Query 参数**:

  | 参数 | 必选 | 类型 | 默认值 | 说明 |
  | --- | --- | --- | --- | --- |
  | releaseKey | 是 | string | - | 客户端当前版本 |
  | timeout | 否 | int | 60 | 超时时间（秒），最大 90 |
  | label | 否 | string | - | 预留 |
  | ip | 否 | string | - | 预留 |

- **请求示例**:
    ```http
    GET /api/v1/envs/pro/apps/app123/clusters/default/namespaces/application/notifications?releaseKey=rel-abc123&timeout=60
    ```

- **响应示例（有更新）**:
    ```json
    {
      "code": 0,
      "message": "success",
      "data": {
        "releaseKey": "rel-def456"
      }
    }
    ```

- **响应示例（超时无更新）**:
    ```json
    {
      "code": 408,
      "message": "timeout",
      "data": null
    }
    ```

---

#### 3. 创建 App
> 创建一个新的应用（对应 app 表）

- **URL**:
  `POST /api/v1/envs/{env}/apps`

- **Path 参数**:

  | 参数 | 必选 | 类型 | 说明 |
  | --- | --- | --- | --- |
  | env | 是 | string | 环境 |

- **Body 参数**:

  | 参数 | 必选 | 类型 | 说明 |
  | --- | --- | --- | --- |
  | appName | 是 | string | 应用名称 |
  | description | 否 | string | 描述 |

- **Body 示例**:
    ```json
    {
      "app_name": "app123",
      "description": "用户服务"
    }
    ```

- **响应示例**:
    ```json
    {
      "code": 0,
      "message": "success",
      "data": {
        "app_name": "app123",
        "description": "处理用户注册登录",
        "create_time": "2025-08-01T14:00:00Z",
        "update_time": "2025-08-01T14:00:00Z"
      }
    }
    ```

---

#### 4、获取App列表

- **URL**:
`GET /api/v1/envs/{env}/apps`

- **HTTP 方法**:
`GET`

- **Path 参数**:
  
  | 参数 | 必选 | 类型   | 说明         |
  | --- | --- | ------ | ------------ |
  | env | 是   | string | 环境名称（如 "dev", "pro"） |

- **Query 参数**:
  
  | 参数   | 必选 | 类型   | 默认值 | 说明                       |
  | ------ | ---- | ------ | ------- | -------------------------- |
  | page   | 否   | int    | 1       | 请求的页码                 |
  | size   | 否   | int    | 20      | 每页显示的应用数量，范围 1-100 |

- **请求示例**:
```bash
curl -X GET "http://localhost:3000/api/v1/envs/dev/apps?page=2&size=15" \
     -H "accept: application/json"
```

- **响应示例**:
成功时返回 HTTP 状态码 `200 OK` 及以下 JSON 响应体：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env": "dev",
    "apps": [
      {
        "app_name": "app123",
        "description": "用户服务",
        "create_time": "2025-08-01T14:00:00Z",
        "update_time": "2025-08-01T14:00:00Z"
      },
      {
        "app_name": "anotherApp",
        "description": "另一个应用的服务",
        "create_time": "2025-07-29T16:30:00Z",
        "update_time": "2025-07-30T09:15:00Z"
      }
    ],
    "total": 2,
    "page": 2,
    "size": 15
  }
}
```

---

#### 5、获取单个App详情

- **URL**:
  `GET /api/v1/envs/{env}/apps/{app_name}`

- **HTTP 方法**:
  `GET`

- **Path 参数**:

  | 参数      | 必选 | 类型   | 说明                                 |
    | --------- | ---- | ------ | ------------------------------------ |
  | `env`     | 是   | string | 环境名称（如 `"dev"`, `"pro"`），自动转为小写，需符合环境命名规范 |
  | `app_name` | 是   | string | 应用名称，需符合应用命名规范（如字母、数字、下划线、短横线等） |

- **Query 参数**:

  无

- **请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/envs/dev/apps/user_service_v1xO" \
     -H "accept: application/json"
```

- **响应示例**:

成功时返回 HTTP 状态码 `200 OK` 及以下 JSON 响应体：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env": "dev",
    "app": {
      "app_name": "user_service_v1xO",
      "description": "合法特殊字符测试",
      "create_time": "2025-08-03T18:54:16+08:00",
      "update_time": "2025-08-03T18:54:16+08:00"
    }
  }
}
```
---

#### 5. 更新App信息

- **URL**:
  `PUT /api/v1/envs/{env}/apps/{app_name}`

- **HTTP 方法**:
  `PUT`

- **Path 参数**:

  | 参数      | 必选 | 类型   | 说明                                 |
    | --------- | ---- | ------ | ------------------------------------ |
  | `env`     | 是   | string | 环境名称（如 `"dev"`, `"pro"`），自动转为小写，需符合环境命名规范 |
  | `app_name` | 是   | string | 应用名称，需符合应用命名规范（如字母、数字、下划线、短横线等） |

- **Body 参数**:

  | 参数          | 必选 | 类型   | 说明                               |
    | ------------- | ---- | ------ | ---------------------------------- |
  | `description` | 否   | string | 应用描述                           |


- **请求示例**:
```bash
curl -X PUT "http://localhost:8080/api/v1/envs/dev/apps/user_service_v1xO" \
     -H "accept: application/json" \
     -H "Content-Type: application/json" \
     -d '{
           "description": "更新后的用户服务",
         }'
```

- **响应示例**:

成功时返回 HTTP 状态码 `200 OK` 及以下 JSON 响应体：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env": "dev",
    "app": {
      "app_name": "user_service_v1xO",
      "description": "更新后的用户服务",
      "create_time": "2025-08-03T18:54:16+08:00",
      "update_time": "2025-08-03T22:00:00+08:00"
    }
  }
}
```


---
### 6. 创建集群（CreateCluster）

- **URL**:  
  `POST /api/v1/envs/{env}/apps/{app_name}/clusters`

- **HTTP 方法**:  
  `POST`

- **Path 参数**:
  
  | 参数       | 必选 | 类型   | 说明 |
  |-----------|------|------|------|
  | `env`      | 是   | string | 环境名（如 `"dev"`, `"prod"`） |
  | `app_name` | 是   | string | 应用名（如 `"slimstor"`） |

- **Body 参数**:
  
  | 参数         | 必选 | 类型   | 说明 |
  |-------------|------|------|------|
  | `cluster_name` | 是   | string | 集群名称（需符合命名规范） |
  | `description`  | 否   | string | 描述信息 |

- **请求示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters" \
     -H "accept: application/json" \
     -H "Content-Type: application/json" \
     -d '{
           "cluster_name": "cluster-1",
           "description": "开发环境主集群"
         }'
```

- **响应示例**（成功）:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env": "dev",
    "cluster": {
      "app_name": "slimstor",
      "cluster_name": "cluster-1",
      "description": "开发环境主集群",
      "create_time": "2025-08-04T10:00:00+08:00",
      "update_time": "2025-08-04T10:00:00+08:00"
    }
  }
}
```

---

### 7. 获取集群列表（ListClusters）

- **URL**:  
  `GET /api/v1/envs/{env}/apps/{app_name}/clusters`

- **HTTP 方法**:  
  `GET`

- **Path 参数**:
  
  | 参数       | 必选 | 类型   | 说明 |
  |-----------|------|------|------|
  | `env`      | 是   | string | 环境名 |
  | `app_name` | 是   | string | 应用名 |

- **Query 参数**:
  
  | 参数   | 必选 | 类型 | 默认值 | 说明 |
  |-------|------|------|--------|------|
  | `page` | 否   | int  | 1      | 页码（≥1） |
  | `size` | 否   | int  | 20     | 每页数量（1-100） |

- **请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters?page=1&size=10" \
     -H "accept: application/json"
```

- **响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env": "dev",
    "clusters": [
      {
        "app_name": "slimstor",
        "cluster_name": "cluster-1",
        "description": "开发环境主集群",
        "create_time": "2025-08-04T10:00:00+08:00",
        "update_time": "2025-08-04T10:00:00+08:00"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  }
}
```

---

### 8. 获取单个集群（GetCluster）

- **URL**:  
  `GET /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}`

- **HTTP 方法**:  
  `GET`

- **Path 参数**:
  
  | 参数           | 必选 | 类型   | 说明 |
  |---------------|------|------|------|
  | `env`          | 是   | string | 环境名 |
  | `app_name`     | 是   | string | 应用名 |
  | `cluster_name` | 是   | string | 集群名 |

- **请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/cluster-1" \
     -H "accept: application/json"
```

- **响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env": "dev",
    "cluster": {
      "app_name": "slimstor",
      "cluster_name": "cluster-1",
      "description": "开发环境主集群",
      "create_time": "2025-08-04T10:00:00+08:00",
      "update_time": "2025-08-04T10:00:00+08:00"
    }
  }
}
```

---

### 9. 删除集群（DeleteCluster）

- **URL**:  
  `DELETE /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}`

- **HTTP 方法**:  
  `DELETE`

- **Path 参数**:
  
  | 参数           | 必选 | 类型   | 说明 |
  |---------------|------|------|------|
  | `env`          | 是   | string | 环境名 |
  | `app_name`     | 是   | string | 应用名 |
  | `cluster_name` | 是   | string | 集群名 |

- **请求示例**:
```bash
curl -X DELETE "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/cluster-1" \
     -H "accept: application/json"
```

- **响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env": "dev",
    "app_name": "slimstor",
    "cluster_name": "cluster-1"
  }
}
```

---
### 10. 创建命名空间（CreateNamespace）

- **URL**:  
  `POST /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces`

- **HTTP 方法**:  
  `POST`

- **Path 参数**:
  
  | 参数           | 必选 | 类型   | 说明 |
  |----------------|------|--------|------|
  | `env`          | 是   | string | 环境名 |
  | `app_name`     | 是   | string | 应用名 |
  | `cluster_name` | 是   | string | 集群名 |

- **请求 Body (JSON)**:
```json
{
  "namespace_name": "default",
  "description": "This is the default namespace."
}
```

- **请求示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/cluster-1/namespaces" \
     -H "accept: application/json" \
     -H "Content-Type: application/json" \
     -d '{
  "namespace_name": "default",
  "description": "This is the default namespace."
}'
```

- **成功响应 (200 OK)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "namespace": {
      "namespace_name": "default",
      "description": "This is the default namespace.",
      "app_name": "slimstor",
      "cluster_name": "cluster-1",
      "create_time": "2025-08-04T16:50:00+08:00",
      "update_time": "2025-08-04T16:50:00+08:00"
    }
  }
}
```

---

### 11. 获取命名空间列表（ListNamespace）

- **URL**:  
  `GET /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces`

- **HTTP 方法**:  
  `GET`

- **Path 参数**:
  
  | 参数           | 必选 | 类型   | 说明 |
  |----------------|------|--------|------|
  | `env`          | 是   | string | 环境名 |
  | `app_name`     | 是   | string | 应用名 |
  | `cluster_name` | 是   | string | 集群名 |

- **Query 参数 (可选)**:
  
  | 参数 | 类型 | 说明 |
  |------|------|------|
  | `page` | int | 分页页数，默认值为1 |
  | `size` | int | 每页显示的记录数，默认值为10 |

- **请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/cluster-1/namespaces?page=1&size=10" \
     -H "accept: application/json"
```

- **成功响应 (200 OK)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "namespaces": [
      {
        "namespace_name": "default",
        "description": "This is the default namespace.",
        "app_name": "slimstor",
        "cluster_name": "cluster-1",
        "create_time": "2025-08-04T16:50:00+08:00",
        "update_time": "2025-08-04T16:50:00+08:00"
      }
    ],
    "total": 100
  }
}
```

---

### 12. 获取单个命名空间（GetNamespace）

- **URL**:  
  `GET /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}`

- **HTTP 方法**:  
  `GET`

- **Path 参数**:
  
  | 参数           | 必选 | 类型   | 说明 |
  |----------------|------|--------|------|
  | `env`          | 是   | string | 环境名 |
  | `app_name`     | 是   | string | 应用名 |
  | `cluster_name` | 是   | string | 集群名 |
  | `namespace_name` | 是 | string | 命名空间名称 |

- **请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/cluster-1/namespaces/default" \
     -H "accept: application/json"
```

- **成功响应 (200 OK)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env" : "dev",
    "app_name": "slimstor",
    "cluster_name": "cluster-1",
    "namespace_name": "default",
    "description": "This is the default namespace.",
    "create_time": "2025-08-04T16:50:00+08:00",
    "update_time": "2025-08-04T16:50:00+08:00"
  }
}
```

---

### 13. 删除命名空间（DeleteNamespace）

- **URL**:  
  `DELETE /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}`

- **HTTP 方法**:  
  `DELETE`

- **Path 参数**:
  
  | 参数           | 必选 | 类型   | 说明 |
  |----------------|------|--------|------|
  | `env`          | 是   | string | 环境名 |
  | `app_name`     | 是   | string | 应用名 |
  | `cluster_name` | 是   | string | 集群名 |
  | `namespace_name` | 是 | string | 命名空间名称 |

- **请求示例**:
```bash
curl -X DELETE "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/cluster-1/namespaces/default" \
     -H "accept: application/json"
```

- **成功响应 (200 OK)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "namespace": {
      "namespace_name": "default",
      "app_name": "slimstor",
      "cluster_name": "cluster-1"
    }
  }
}
```

---

### 14. 创建 Item（SetItem）

- **URL**:  
  `POST /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}/items`

- **HTTP 方法**:  
  `POST`

- **Path 参数**:
  
  | 参数           | 必选 | 类型   | 说明 |
  |----------------|------|--------|------|
  | `env`          | 是   | string | 环境名 |
  | `app_name`     | 是   | string | 应用名 |
  | `cluster_name` | 是   | string | 集群名 |
  | `namespace_name`| 是   | string | 命名空间名称 |

- **请求 Body (JSON)**:
```json
{
  "key": "item_key",
  "value": "This is the value of the item."
}
```

- **请求示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/shenzhen/namespaces/default/items" \
     -H "accept: application/json" \
     -H "Content-Type: application/json" \
     -d '{
  "key": "item_key",
  "value": "This is the value of the item."
}'
```

- **成功响应 (200 OK)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "item": {
      "key": "item_key",
      "value": "This is the value of the item.",
      "namespace_name": "default",
      "create_time": "2025-08-05T00:00:00+08:00",
      "update_time": "2025-08-05T00:00:00+08:00"
    }
  }
}
```
---

## 17. 获取命名空间下所有 Items（ListItem）

**URL:**  
`GET /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}/items`

**HTTP 方法:**  
`GET`

---

### Path 参数

| 参数 | 必选 | 类型 | 说明 |
|------|------|------|------|
| `env` | 是 | `string` | 环境名，例如：`dev`, `prod` |
| `app_name` | 是 | `string` | 应用名，例如：`slimstor` |
| `cluster_name` | 是 | `string` | 集群名，例如：`shenzhen`, `cluster-1` |
| `namespace_name` | 是 | `string` | 命名空间名称，例如：`default`, `database` |

---

### Query 参数（可选）

| 参数 | 必选 | 类型 | 默认值 | 说明 |
|------|------|------|--------|------|
| `page` | 否 | `int` | `1` | 分页页码，从 1 开始 |
| `size` | 否 | `int` | `20` | 每页返回数量，最大支持 `100` |
| `search` | 否 | `string` | `""` | 按 `key` 字段进行模糊搜索（不区分大小写） |

> ⚠️ 说明：`page` 和 `size` 用于分页控制；`search` 支持子串匹配，如 `search=database` 可匹配 `database.host`、`app.database.url` 等。

---

### 请求示例

```bash
curl -X GET "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/shenzhen/namespaces/default/items?page=1&size=10&search=database" \
     -H "accept: application/json"
```

---

### 成功响应 (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "key": "database.host",
        "value": "10.0.0.1",
        "created_at": "2025-01-01T10:00:00Z",
        "updated_at": "2025-01-01T10:00:00Z"
      },
      {
        "key": "database.port",
        "value": "3306",
        "created_at": "2025-01-01T10:00:00Z",
        "updated_at": "2025-01-01T10:00:00Z"
      },
      {
        "key": "database.user",
        "value": "admin",
        "created_at": "2025-01-01T10:00:00Z",
        "updated_at": "2025-01-01T10:00:00Z"
      }
    ],
    "total": 3,
    "page": 1,
    "size": 10
  }
}
```

---

### 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | `int` | 返回码，`0` 表示成功 |
| `message` | `string` | 提示信息，成功时为 `"success"` |
| `data` | `object` | 返回数据对象 |
| &nbsp;&nbsp;`items` | `array` | 配置项列表 |
| &nbsp;&nbsp;&nbsp;&nbsp;`key` | `string` | 配置项的键 |
| &nbsp;&nbsp;&nbsp;&nbsp;`value` | `string` | 配置项的值 |
| &nbsp;&nbsp;&nbsp;&nbsp;`created_at` | `string` (ISO 8601) | 创建时间 |
| &nbsp;&nbsp;&nbsp;&nbsp;`updated_at` | `string` (ISO 8601) | 最后更新时间 |
| &nbsp;&nbsp;`total` | `int` | 总记录数（用于分页） |
| &nbsp;&nbsp;`page` | `int` | 当前页码 |
| &nbsp;&nbsp;`size` | `int` | 每页数量 |

---

### 错误响应示例

#### 404 Not Found - 命名空间不存在
```json
{
  "code": 404,
  "message": "namespace not found",
  "data": null
}
```

#### 400 Bad Request - 参数错误（如 page <= 0）
```json
{
  "code": 400,
  "message": "invalid page or size",
  "data": null
}
```

#### 500 Internal Server Error - 服务端异常
```json
{
  "code": 500,
  "message": "internal server error",
  "data": null
}
```

---
### 16. 获取 Item（GetItem）

- **URL**:  
  `GET /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}/items/{key}`

- **HTTP 方法**:  
  `GET`

- **Path 参数**:

| 参数           | 必选 | 类型   | 说明 |
|----------------|------|--------|------|
| `env`          | 是   | string | 环境名 |
| `app_name`     | 是   | string | 应用名 |
| `cluster_name` | 是   | string | 集群名 |
| `namespace_name`| 是   | string | 命名空间名称 |
| `key`          | 是   | string | Item 的键 |

- **请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/shenzhen/namespaces/default/items/item_key" \
     -H "accept: application/json"
```

- **成功响应 (200 OK)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
      "key": "item_key",
      "value": "This is the value of the item.",
      "namespace_name": "default",
      "create_time": "2025-08-05T00:00:00+08:00",
      "update_time": "2025-08-05T00:00:00+08:00"
  }
}
```

---

### 17. 删除 Item（DeleteItem）

- **URL**:  
  `DELETE /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}/items/{key}`

- **HTTP 方法**:  
  `DELETE`

- **Path 参数**:

| 参数           | 必选 | 类型   | 说明 |
|----------------|------|--------|------|
| `env`          | 是   | string | 环境名 |
| `app_name`     | 是   | string | 应用名 |
| `cluster_name` | 是   | string | 集群名 |
| `namespace_name`| 是   | string | 命名空间名称 |
| `key`          | 是   | string | Item 的键 |

- **请求示例**:
```bash
curl -X DELETE "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/cluster-1/namespaces/default/items/item_key" \
     -H "accept: application/json"
```

- **成功响应 (200 OK)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "key": "item_key",
    "namespace_name": "default"
  }
}
```
---

### 18. 发布命名空间配置（PublishRelease）

- **URL**:  
  `POST /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}/releases`

- **HTTP 方法**:  
  `POST`

- **Path 参数**:

  | 参数             | 必选 | 类型   | 说明 |
    |------------------|------|--------|------|
  | `env`            | 是   | string | 环境名，如 `dev`、`prod` |
  | `app_name`       | 是   | string | 应用名称 |
  | `cluster_name`   | 是   | string | 集群名称 |
  | `namespace_name` | 是   | string | 命名空间名称 |

- **请求 Body (JSON)**:
```json
{
  "operator": "stevenrao",
  "release_name": "slimstor.shenzhen.default.20250805120000",
  "comment": "This is a 备注."
}
```

| 字段         | 必选 | 类型   | 说明 |
  |--------------|------|--------|------|
| `operator`     | 是   | string | 操作人，用于审计和展示 |
| `release_name` | 否   | string | 自定义发布名称，若不传可由服务端生成 |
| `comment`      | 否   | string | 发布备注，支持中文、特殊字符等 |

- **请求示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/shenzhen/namespaces/default/releases" \
     -H "accept: application/json" \
     -H "Content-Type: application/json" \
     -d '{
  "operator": "stevenrao",
  "release_name": "slimstor.shenzhen.default.20250805120000",
  "comment": "This is a 备注."
}'
```

- **成功响应 (200 OK)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env": "dev",
    "app_name": "slimstor",
    "cluster_name": "shenzhen",
    "namespace_name": "default",
    "release_name": "slimstor.shenzhen.default.20250805120000",
    "release_id": "r-20250805120000abcdef123456",
    "release_time": "2025-08-05T12:00:00+08:00"
  }
}
```
---

### 19. 列出命名空间的所有发布版本（ListReleases）

- **URL**:  
  `GET /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}/releases`

- **HTTP 方法**:  
  `GET`

- **Path 参数**:

  | 参数             | 必选 | 类型   | 说明                          |
    |------------------|------|--------|-------------------------------|
  | `env`            | 是   | string | 环境名，如 `dev`、`prod`      |
  | `app_name`       | 是   | string | 应用名称                      |
  | `cluster_name`   | 是   | string | 集群名称                      |
  | `namespace_name` | 是   | string | 命名空间名称                  |

- **Query 参数 (可选)**:

  | 参数         | 必选 | 类型   | 说明                           |
    |--------------|------|--------|--------------------------------|
  | `page`       | 否   | int    | 分页查询的页码，默认为 1        |
  | `size`       | 否   | int    | 每页显示的数据条数，默认为 20  |
  | `sort`       | 否   | string | 排序字段及顺序，例如 `release_time,desc` |

- **请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/shenzhen/namespaces/default/releases?page=1&size=20&sort=release_time,desc" \
     -H "accept: application/json"
```

- **成功响应 (200 OK)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "env": "dev",
    "size": 100,
    "total": 14,
    "page": 1,
    "releases": [
      {
        "app_name": "slimstor",
        "cluster_name": "shenzhen",
        "namespace_name": "default",
        "release_id": "01987ac7-cbaf-7ce1-8168-99ca8b11d360",
        "release_name": "slimstor.shenzhen.default.2025080512000022",
        "release_time": "2025-08-05T15:09:31+08:00",
        "operator": "stevenrao",
        "comment": "再发布一次."
      }
    ]
  }
}
```

---

### 19. 获取发布详情（GetRelease）

- **URL**:  
  `GET /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}/releases/{release_id}`

- **HTTP 方法**:  
  `GET`

- **Path 参数**:

  | 参数             | 必选 | 类型   | 说明 |
    |------------------|------|--------|------|
  | `env`            | 是   | string | 环境名，如 `dev`、`prod` |
  | `app_name`       | 是   | string | 应用名称 |
  | `cluster_name`   | 是   | string | 集群名称 |
  | `namespace_name` | 是   | string | 命名空间名称 |
  | `release_id`     | 是   | string | 发布 ID，用于比对增量内容；传空字符串或无效 ID 时返回完整配置 |

- **Query 参数**:  
  无

- **请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/shenzhen/namespaces/default/releases/rel-20250805120000" \
     -H "accept: application/json"
```

- **成功响应 (HTTP 200)**:

  返回最新发布版本的完整配置及与指定 `release_id` 的增量差异（如提供有效历史版本）。

```json
{
  "env": "dev",
  "app_name": "slimstor",
  "cluster_name": "shenzhen",
  "namespace_name": "default",
  "release_id": "rel-20250805120001",
  "release_name": "slimstor.shenzhen.default.20250805120001",
  "release_time": "2025-08-05T12:00:01Z",
  "operator": "stevenrao",
  "comment": "This is a 备注.",
  "items": [
    {
      "key": "db.host",
      "value": "10.10.1.1"
    },
    {
      "key": "db.port",
      "value": "3306"
    }
  ],
  "changed": {
    "added": ["db.port"],
    "updated": ["db.host"],
    "deleted": ["db.user"]
  }
}
```

---

### 20. 回滚发布（RollbackRelease）

- **URL**:  
  `POST /api/v1/envs/{env}/apps/{app_name}/clusters/{cluster_name}/namespaces/{namespace_name}/releases/{release_id}/rollback`

- **HTTP 方法**:  
  `POST`

- **Path 参数**:

  | 参数             | 必选 | 类型   | 说明 |
    |------------------|------|--------|------|
  | `env`            | 是   | string | 环境名，如 `dev`、`prod` |
  | `app_name`       | 是   | string | 应用名称 |
  | `cluster_name`   | 是   | string | 集群名称 |
  | `namespace_name` | 是   | string | 命名空间名称 |
  | `release_id`     | 是   | string | 待回滚到的目标发布版本 ID |

- **请求 Body (JSON)**:
```json
{
  "operator": "stevenrao",
  "comment": "回滚到稳定版本 v1.2.0"
}
```

| 字段       | 必选 | 类型   | 说明 |
|------------|------|--------|------|
| `operator` | 是   | string | 操作人，用于审计和记录 |
| `comment`  | 否   | string | 回滚原因或备注，可选字段，支持中文和特殊字符 |

- **请求示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/envs/prod/apps/slimstor/clusters/shenzhen/namespaces/default/releases/rel-20250805100000/rollback" \
     -H "accept: application/json" \
     -H "Content-Type: application/json" \
     -d '{
  "operator": "stevenrao",
  "comment": "回滚到稳定版本 v1.2.0"
}'
```

- **成功响应 (HTTP 200)**:

  回滚成功后，返回新生成的发布版本信息。

```json
{
  "env": "prod",
  "app_name": "slimstor",
  "cluster_name": "shenzhen",
  "namespace_name": "default",
  "release_id": "rel-20250805153000",
  "release_name": "rollback_to_rel-20250805100000",
  "release_time": "2025-08-05T15:30:00Z",
  "operator": "stevenrao",
  "comment": "回滚到稳定版本 v1.2.0"
}
```
---

---


