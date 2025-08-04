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
curl -X POST "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/cluster-1/namespaces/default/items" \
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

### 15. 获取 Item（GetItem）

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
curl -X GET "http://localhost:8080/api/v1/envs/dev/apps/slimstor/clusters/cluster-1/namespaces/default/items/item_key" \
     -H "accept: application/json"
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

### 16. 删除 Item（DeleteItem）

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
    "item": {
      "key": "item_key",
      "namespace_name": "default"
    }
  }
}
```

---

#### 9. 发布 Namespace
> 将当前 item 表中的配置发布为一个新版本（对应 namespace_release）

- **URL**:
  `POST /api/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/releases`

- **Path 参数**:

  | 参数 | 必选 | 类型 | 说明 |
  | --- | --- | --- | --- |
  | env | 是 | string | 环境 |
  | appId | 是 | string | 应用ID |
  | clusterName | 是 | string | 集群名称 |
  | namespaceName | 是 | string | 命名空间名称 |

- **Body 参数**:

  | 参数 | 必选 | 类型 | 说明 |
  | --- | --- | --- | --- |
  | operator | 是 | string | 操作人 |
  | comment | 否 | string | 备注 |

- **Body 示例**:
    ```json
    {
      "operator": "zhangsan",
      "comment": "发布数据库连接配置"
    }
    ```

- **响应示例**:
    ```json
    {
      "code": 0,
      "message": "released",
      "data": {
        "releaseKey": "20250801-release-xyz",
        "namespaceName": "application",
        "releaseTime": "2025-08-01T15:20:00Z"
      }
    }
    ```

---

#### 10. 读取单个 Key-Value

- **URL**:
  `GET /api/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items/{key}`

- **Path 参数**:

  | 参数 | 必选 | 类型 | 说明 |
  | --- | --- | --- | --- |
  | env | 是 | string | 环境 |
  | appId | 是 | string | 应用ID |
  | clusterName | 是 | string | 集群名称 |
  | namespaceName | 是 | string | 命名空间名称 |
  | key | 是 | string | 配置项Key |

- **请求示例**:
    ```http
    GET /api/v1/envs/pro/apps/app123/clusters/default/namespaces/application/items/db.url
    ```

- **响应示例**:
    ```json
    {
      "code": 0,
      "message": "success",
      "data": {
        "key": "db.url",
        "value": "jdbc:mysql://pro"
      }
    }
    ```

---

#### 11. 写入单个 Key-Value
> ⚠️ 仅更新 item 表，不触发发布

- **URL**:
  `PUT /api/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items/{key}`

- **Path 参数**:

  | 参数 | 必选 | 类型 | 说明 |
  | --- | --- | --- | --- |
  | env | 是 | string | 环境 |
  | appId | 是 | string | 应用ID |
  | clusterName | 是 | string | 集群名称 |
  | namespaceName | 是 | string | 命名空间名称 |
  | key | 是 | string | 配置项Key |

- **Body 参数**:

  | 参数 | 必选 | 类型 | 说明 |
  | --- | --- | --- | --- |
  | value | 是 | string | 新值 |
  | operator | 是 | string | 操作人 |

- **Body 示例**:
    ```json
    {
      "value": "new_value",
      "operator": "lisi"
    }
    ```

- **响应示例**:
    ```json
    {
      "code": 0,
      "message": "item updated",
      "data": {
        "key": "db.url",
        "value": "new_value",
        "kvId": 999888777
      }
    }
    ```

---

#### 12. 删除单个 Key-Value
> ⚠️ 仅删除 item 表中的条目，不触发发布

- **URL**:
  `DELETE /api/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items/{key}`

- **Path 参数**:

  | 参数 | 必选 | 类型 | 说明 |
  | --- | --- | --- | --- |
  | env | 是 | string | 环境 |
  | appId | 是 | string | 应用ID |
  | clusterName | 是 | string | 集群名称 |
  | namespaceName | 是 | string | 命名空间名称 |
  | key | 是 | string | 配置项Key |

- **请求示例**:
    ```http
    DELETE /api/v1/envs/pro/apps/app123/clusters/default/namespaces/application/items/debug.flag
    ```

- **响应示例**:
    ```json
    {
      "code": 0,
      "message": "item deleted",
      "data": null
    }
    ```

---


