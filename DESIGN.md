# 软件工程设计文档

## 目标

- 设计一个基于 Go 语言的电商后端服务系统，提供商品管理、用户管理、购物车和订单管理的完整功能
- 采用清晰的分层架构（Handler-Service-Repository），实现业务逻辑与数据访问的分离
- 提供 RESTful API 接口，支持前端或移动端的调用
- 实现基于 JWT 的用户认证和授权机制，保障系统安全性
- 支持多种数据库（MySQL、PostgreSQL），提供灵活的数据存储方案
- 构建可扩展、可维护的系统架构，为未来业务迭代提供基础

## 非目标

- 本系统不提供前端界面，仅提供 API 接口
- 本系统不涉及支付网关的集成（可在后续迭代中添加）
- 本系统不提供分布式事务支持（当前为单体应用）
- 本系统不支持多地域部署和数据同步（当前为单节点部署）
- 本系统不提供实时消息推送功能

## 背景

### 目标读者

本文档的目标读者包括：
- 项目开发团队成员
- 系统架构师和技术评审人员
- 运维和部署人员

### 项目背景

随着电商业务的快速发展，需要一个稳定、可扩展的后端服务来支撑商品管理、用户管理、购物车和订单等核心业务功能。本项目旨在构建一个现代化的、基于微服务理念的后端系统，具有以下特点：

1. **清晰的架构设计**：采用经典的三层架构（Handler-Service-Repository），职责分明，便于维护和扩展
2. **现代化的技术栈**：使用 Go 语言、Gin 框架、GORM ORM、JWT 认证等成熟技术
3. **依赖注入模式**：使用 Uber dig 库实现依赖注入，降低模块间耦合
4. **标准化的响应格式**：统一的 API 响应结构，便于前端处理

### 项目必要性

- **业务需求**：需要一个可靠的后端系统来支撑电商核心业务流程
- **技术选型**：Go 语言在高并发场景下具有优异的性能表现
- **架构升级**：采用清晰的分层架构，为未来的微服务化改造打下基础
- **可维护性**：通过依赖注入和清晰的代码组织，提高系统的可维护性

## 总体设计

### 系统架构概览

本系统采用经典的三层架构模式，从上到下依次为：

1. **Handler 层（控制器层）**
   - 负责接收 HTTP 请求，参数验证和响应格式化
   - 调用 Service 层完成业务逻辑
   - 返回统一格式的 JSON 响应

2. **Service 层（业务逻辑层）**
   - 实现核心业务逻辑，如用户注册/登录、商品 CRUD、订单处理等
   - 调用 Repository 层完成数据持久化
   - 处理业务规则和数据转换

3. **Repository 层（数据访问层）**
   - 封装数据库操作，提供 CRUD 接口
   - 使用 GORM 进行 ORM 映射
   - 屏蔽具体的数据库实现细节

### 系统上下文关系图

```
┌─────────────────────────────────────────────────────────────┐
│                         Client                               │
│                    (Web/Mobile App)                          │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTPS/HTTP
                           │ RESTful API
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                     Gin Web Server                           │
│  ┌─────────────┐    ┌──────────────┐    ┌──────────────┐   │
│  │   Router    │───▶│  Middleware  │───▶│   Handler    │   │
│  │             │    │  - Auth      │    │  - User      │   │
│  │             │    │  - Logger    │    │  - Commodity │   │
│  └─────────────┘    │  - Recovery  │    │  - Cart      │   │
│                     └──────────────┘    │  - Order     │   │
│                                         └──────┬───────┘   │
│                                                │            │
│                                         ┌──────▼───────┐   │
│                                         │   Service    │   │
│                                         │  - User      │   │
│                                         │  - Commodity │   │
│                                         │  - Cart      │   │
│                                         │  - Order     │   │
│                                         └──────┬───────┘   │
│                                                │            │
│                                         ┌──────▼───────┐   │
│                                         │  Repository  │   │
│                                         │  - User      │   │
│                                         │  - Commodity │   │
│                                         │  - Cart      │   │
│                                         │  - Order     │   │
│                                         └──────┬───────┘   │
└────────────────────────────────────────────────┼───────────┘
                                                 │ GORM
                                                 ▼
┌─────────────────────────────────────────────────────────────┐
│                   Database Layer                             │
│  ┌──────────────────┐         ┌──────────────────┐          │
│  │      MySQL       │   or    │    PostgreSQL    │          │
│  └──────────────────┘         └──────────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

### 核心技术栈

- **Web 框架**：Gin - 高性能的 Go Web 框架
- **ORM**：GORM - 功能完善的 Go ORM 库
- **依赖注入**：Uber dig - 依赖注入容器
- **认证**：JWT (JSON Web Token) - 无状态的用户认证
- **日志**：Logrus - 结构化日志库
- **配置管理**：Viper - 灵活的配置管理库

### 核心设计理念

1. **依赖注入**：通过 dig 容器管理所有组件的生命周期，实现松耦合
2. **单一职责**：每个模块只负责特定的功能，便于测试和维护
3. **统一响应**：所有 API 返回统一的响应格式（code, message, data）
4. **中间件机制**：通过中间件实现横切关注点（认证、日志、错误处理）
5. **配置外部化**：数据库连接、服务端口等配置通过 YAML 文件管理

## 详细设计

### 各子模块的设计

#### 1. 用户模块 (User Module)

**功能职责**：
- 用户注册：接收账号、密码、姓名，创建新用户
- 用户登录：验证账号密码，生成 JWT Token
- 密码加密：使用 bcrypt 对密码进行哈希存储

**核心流程**：

用户注册流程：
```
Client → Handler.Register → Service.Register → Repository.CreateUser → Database
                                ↓
                        bcrypt加密密码
```

用户登录流程：
```
Client → Handler.Login → Service.Login → Repository.FindUserByAccount
                              ↓
                      验证密码 + 生成JWT Token
                              ↓
                        返回Token给客户端
```

**关键实现**：
- 密码使用 bcrypt 加密，成本因子为 10
- JWT Token 包含用户 ID、账号信息，当前代码配置为 999 小时（**⚠️ 安全风险：该配置仅用于记录现有实现，不代表推荐做法。强烈建议即使在开发环境也缩短为 24-48 小时，生产环境必须缩短为 1-24 小时，并实现 Token 刷新机制**）
- JWT 签名密钥：当前代码使用固定密钥 `gee`（**⚠️ 严重安全风险：该配置仅用于记录现有实现，不代表推荐做法。生产环境必须使用强随机密钥（至少 256 位），并通过环境变量配置，切勿将密钥硬编码在代码中。建议开发环境也使用环境变量管理密钥**）

#### 2. 商品模块 (Commodity Module)

**功能职责**：
- 创建商品：添加新商品信息（名称、价格、库存）
- 更新商品：修改商品信息
- 删除商品：软删除或硬删除商品
- 查询商品：按名称查询、列表查询

**数据模型**：
```go
type Commodity struct {
    ID        int       // 商品ID
    Name      string    // 商品名称
    Price     float64   // 商品价格
    Stock     int       // 库存数量
    Status    bool      // 商品状态（上架/下架）
    CreatedAt time.Time // 创建时间
    UpdateAt  time.Time // 更新时间
}
```

**关键特性**：
- 支持商品状态管理（上架/下架）
- 库存管理，防止超卖
- 软删除支持，保留历史数据

#### 3. 购物车模块 (Cart Module)

**功能职责**：
- 添加商品到购物车
- 移除购物车中的商品
- 更新购物车商品数量
- 查询用户购物车

**数据模型**：
```go
type Cart struct {
    Id          int       // 购物车项ID
    UserId      int       // 用户ID
    CommodityId int       // 商品ID
    Quantity    int       // 数量
    CreatedAt   time.Time // 创建时间
    UpdateAt    time.Time // 更新时间
}
```

**业务规则**：
- 每个用户可以有多个购物车项
- 同一商品只能有一个购物车项，通过数量控制
- 购物车与用户关联，需要登录后才能操作

#### 4. 订单模块 (Order Module)

**功能职责**：
- 创建订单：从购物车生成订单
- 更新订单状态：待支付、已支付、已发货、已完成、已取消
- 删除订单：取消订单
- 查询订单：按用户、状态等条件查询

**数据模型**：
```go
type Order struct {
    Id          int       // 订单ID
    UserId      int       // 用户ID
    CommodityId int       // 商品ID
    Quantity    int       // 数量
    TotalPrice  float64   // 总价
    Address     string    // 收货地址
    Status      string    // 订单状态
    CreatedAt   time.Time // 创建时间
    UpdateAt    time.Time // 更新时间
}
```

**订单状态流转**：
```
待支付 → 已支付 → 已发货 → 已完成
   ↓
已取消
```

#### 5. 认证中间件 (Auth Middleware)

**功能职责**：
- 验证请求头中的 JWT Token
- 解析 Token 获取用户信息
- 将用户信息注入到上下文中

**工作流程**：
```
Request → 提取Authorization头 → 验证Bearer Token → 解析JWT
                                      ↓
                            验证签名 + 检查过期时间
                                      ↓
                        将用户ID和账号存入Context → Next Handler
```

**安全措施**：
- Token 必须以 "Bearer " 开头
- 验证 Token 签名和过期时间
- 失败时返回 401 Unauthorized

### API 接口设计

#### 公开接口（无需认证）

| 方法 | 路径 | 功能 | 请求体 | 响应 |
|------|------|------|--------|------|
| POST | /v1/register | 用户注册 | `{account, password, name}` | `{code, message, data}` |
| POST | /v1/login | 用户登录 | `{account, password}` | `{code, message, data: {token}}` |

#### 认证接口（需要 JWT Token）

**商品相关**：
| 方法 | 路径 | 功能 | 请求体 | 响应 |
|------|------|------|--------|------|
| POST | /v1/createCommodity | 创建商品 | `{name, price, stock}` | `{code, message, data}` |
| POST | /v1/updateCommodity | 更新商品 | `{id, name, price, stock}` | `{code, message, data}` |
| GET | /v1/listCommodity | 商品列表 | - | `{code, message, data: []}` |
| DELETE | /v1/deleteCommodity | 删除商品 | `{id}` | `{code, message, data}` |
| GET | /v1/getCommodity | 查询商品 | `?name=xxx` | `{code, message, data}` |

**购物车相关**：
| 方法 | 路径 | 功能 | 请求体 | 响应 |
|------|------|------|--------|------|
| POST | /v1/addToCart | 添加到购物车 | `{commodityId, quantity}` | `{code, message, data}` |
| DELETE | /v1/removeFromCart | 移除购物车 | `{cartId}` | `{code, message, data}` |
| PUT | /v1/updateCart | 更新购物车 | `{cartId, quantity}` | `{code, message, data}` |
| GET | /v1/getCart | 查询购物车 | - | `{code, message, data: []}` |

**订单相关**：
| 方法 | 路径 | 功能 | 请求体 | 响应 |
|------|------|------|--------|------|
| POST | /v1/createOrder | 创建订单 | `{commodityId, quantity, address}` | `{code, message, data}` |
| PUT | /v1/updateOrder | 更新订单状态 | `{orderId, status}` | `{code, message, data}` |
| DELETE | /v1/deleteOrder | 删除订单 | `{orderId}` | `{code, message, data}` |
| GET | /v1/getOrder | 查询订单 | `?userId=xxx` | `{code, message, data: []}` |

#### 统一响应格式

```json
{
  "code": 0,           // 业务错误码，0表示成功
  "message": "success", // 提示信息
  "data": {}           // 响应数据
}
```

**错误码设计**：
- 格式：模块(2位) + 类型(2位) + 序号(2位)
- 通用错误码：0xxxxx（如：100000 内部错误，100003 未授权）
- 用户模块：10xxxx（如：101001 用户不存在，101003 密码错误）
- 商品模块：20xxxx（如：201001 商品不存在，201002 创建失败）

### 存储设计

#### 数据库选型

支持多种关系型数据库：
- **MySQL**：生产环境推荐，成熟稳定
- **PostgreSQL**：支持更丰富的数据类型和查询功能

#### 数据模型设计

**用户表 (users)**：
```sql
CREATE TABLE users (
    uid INT PRIMARY KEY AUTO_INCREMENT,
    account VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**商品表 (commodities)**：
```sql
CREATE TABLE commodities (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

**购物车表 (carts)**：
```sql
CREATE TABLE carts (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    commodity_id INT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(uid),
    FOREIGN KEY (commodity_id) REFERENCES commodities(id),
    UNIQUE KEY uk_user_commodity (user_id, commodity_id)
);
```

**订单表 (orders)**：
```sql
CREATE TABLE orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    commodity_id INT NOT NULL,
    quantity INT NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    address VARCHAR(200),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(uid),
    FOREIGN KEY (commodity_id) REFERENCES commodities(id)
);
```

#### 数据规模预估

**初期（前 6 个月）**：
- 用户数：1000 - 10000
- 商品数：100 - 1000
- 订单数：1000 - 10000
- 购物车项：5000 - 50000

**中期（1-2 年）**：
- 用户数：10 万
- 商品数：1 万
- 订单数：100 万
- 购物车项：50 万

**扩展性考虑**：
- 目前为单表设计，满足中小规模业务
- 当订单量达到千万级别时，需要考虑分库分表
- 分表策略：按用户 ID 或订单创建时间进行 Hash 分片
- 读写分离：主库写入，从库读取，提升查询性能

#### 索引设计

**用户表**：
- 主键索引：uid
- 唯一索引：account（用于登录查询）

**商品表**：
- 主键索引：id
- 普通索引：name（用于名称搜索）
- 普通索引：status（用于状态过滤）

**购物车表**：
- 主键索引：id
- 唯一索引：(user_id, commodity_id)（防止重复添加）
- 外键索引：user_id, commodity_id

**订单表**：
- 主键索引：id
- 普通索引：user_id（用户订单查询）
- 普通索引：status（订单状态查询）
- 复合索引：(user_id, created_at)（用户订单时间查询）

### 考虑过的其他设计方案

#### 1. 认证方案选择

**方案 A：Session + Cookie**
- 优点：服务端可控，可主动失效
- 缺点：需要 Redis 等缓存存储 Session，增加系统复杂度

**方案 B：JWT Token（已选择）**
- 优点：无状态，易于水平扩展，无需额外存储
- 缺点：无法主动失效，Token 相对较大

**选择原因**：JWT 更适合 RESTful API，无需额外的存储依赖，便于扩展

#### 2. 依赖注入方案

**方案 A：手动构造依赖链**
- 优点：简单直观，无额外依赖
- 缺点：随着模块增多，main 函数会变得臃肿

**方案 B：Uber dig（已选择）**
- 优点：自动解析依赖，代码简洁，易于测试
- 缺点：增加学习成本

**选择原因**：dig 提供了优雅的依赖管理，便于单元测试和模块替换

#### 3. 数据库访问方案

**方案 A：原生 SQL**
- 优点：性能最优，控制精确
- 缺点：开发效率低，易出错

**方案 B：GORM（已选择）**
- 优点：开发效率高，自动迁移，类型安全
- 缺点：复杂查询性能略差

**选择原因**：GORM 在开发效率和性能之间取得了良好平衡，适合快速迭代

## 系统 SLA

### 准确性 (Accuracy)

- **目标错误率**：< 0.01%
- **数据一致性**：确保订单创建、库存扣减的事务一致性
- **业务逻辑正确性**：通过单元测试和集成测试保证

### 系统容量 (Capacity)

- **当前支持**：500 QPS（单节点）
- **峰值支持**：1000 QPS（通过负载均衡可水平扩展）
- **并发连接**：5000 个并发连接
- **数据库连接池**：最大 100 个连接

**扩展方案**：
1. 通过负载均衡部署多个应用实例
2. 数据库读写分离，主从复制
3. 引入缓存层（Redis）减轻数据库压力

### 延迟 (Latency)

- **P50 延迟**：< 50ms（简单查询）
- **P99 延迟**：< 200ms（复杂查询）
- **P999 延迟**：< 500ms
- **登录接口**：< 100ms（包含密码验证）

**优化措施**：
- 数据库查询优化（索引、查询语句）
- 响应数据压缩（Gzip）
- 使用连接池减少数据库连接开销

### 可用性 (Availability)

- **目标可用性**：99.9%（即每月停机时间 < 43.2 分钟）
- **恢复时间目标 (RTO)**：< 5 分钟
- **恢复点目标 (RPO)**：< 10 分钟

**保障措施**：
1. 健康检查接口：`/health`
2. 应用自动重启机制
3. 数据库主从切换
4. 定时数据备份

### SLA 验证方案

#### 压力测试

使用工具：Apache Bench (ab) 或 wrk

测试场景：
1. **用户注册/登录**：1000 并发，持续 1 分钟
2. **商品查询**：500 并发，持续 5 分钟
3. **创建订单**：200 并发，持续 2 分钟

期望结果：
- 成功率 > 99.99%
- P99 延迟符合 SLA 要求
- 无内存泄漏、连接泄漏

#### 监控指标

- **QPS/TPS**：每秒请求数/事务数
- **响应时间分布**：P50, P90, P99, P999
- **错误率**：4xx, 5xx 错误占比
- **系统资源**：CPU、内存、网络、磁盘 I/O

## 系统依赖

| 系统依赖 | 说明 | 如果依赖不可用或性能下降的影响 | 处理方案 |
|---------|------|---------------------------|---------|
| **MySQL/PostgreSQL 数据库** | 核心数据存储 | 系统完全不可用，无法读写数据 | 1. 主从复制，自动故障切换<br>2. 定时备份，快速恢复<br>3. 数据库连接池重试机制 |
| **Gin Web 框架** | HTTP 服务器框架 | 应用无法启动 | 无需处理，属于应用核心依赖 |
| **GORM** | ORM 框架 | 数据库操作失败 | 代码层面处理数据库错误，返回友好提示 |
| **JWT 库** | Token 生成和验证 | 无法认证用户 | 1. 使用稳定版本<br>2. 错误处理和日志记录 |
| **Logrus** | 日志库 | 日志无法记录（不影响业务） | 降级到标准输出 |
| **Viper** | 配置管理 | 应用无法启动（缺少配置） | 1. 提供默认配置<br>2. 配置文件校验机制 |

### 外部依赖管理

**依赖版本管理**：
- 使用 Go Modules 管理依赖版本
- 定期更新依赖包，修复安全漏洞
- 使用 `go.sum` 确保依赖包的完整性

**依赖隔离**：
- 通过接口抽象外部依赖
- 便于测试时 Mock 外部依赖
- 降低对具体实现的耦合

## 监控和告警

### 指标上报

#### 系统指标

- **CPU 使用率**：监控应用 CPU 占用，阈值 70%
- **内存使用率**：监控内存占用和 GC 情况，阈值 80%
- **Goroutine 数量**：监控协程泄漏，异常增长时告警
- **数据库连接数**：监控连接池使用情况

#### 业务指标

- **API 调用次数**：按接口统计 QPS
- **API 响应时间**：P50, P90, P99 延迟
- **API 错误率**：按错误码统计
- **用户注册/登录数**：业务增长指标
- **订单创建数**：核心业务指标

#### 日志上报

使用 Logrus 记录以下日志：
- **Info 级别**：用户登录、订单创建等关键操作
- **Error 级别**：数据库错误、业务异常
- **Debug 级别**：请求参数、响应数据（仅开发环境）

日志格式：
```
[INFO] 2024.10.15 08:30:25 user admin try to log in
[INFO] 2024.10.15 08:30:25 user login success: admin
[ERROR] 2024.10.15 08:30:26 database error: connection timeout
```

### 监控方案

#### 推荐工具

1. **Prometheus + Grafana**
   - Prometheus 采集指标
   - Grafana 可视化展示
   - 配置告警规则

2. **ELK Stack（Elasticsearch + Logstash + Kibana）**
   - 集中日志管理
   - 日志搜索和分析
   - 异常日志告警

#### 监控看板

**系统健康看板**：
- CPU、内存、网络实时监控
- 数据库连接数、慢查询统计
- Goroutine 数量变化趋势

**业务监控看板**：
- 接口 QPS 和响应时间
- 错误率和错误码分布
- 用户活跃度和订单量

**使用对象**：
- 系统健康看板：运维团队实时监控
- 业务监控看板：产品和开发团队查看

### 告警策略

#### 告警规则

| 指标 | 阈值 | 级别 | 通知对象 |
|-----|------|------|---------|
| CPU 使用率 | > 80% 持续 5 分钟 | 警告 | 运维团队 |
| 内存使用率 | > 90% | 严重 | 运维团队 |
| API 错误率 | > 1% | 严重 | 开发团队 |
| API P99 延迟 | > 500ms | 警告 | 开发团队 |
| 数据库连接失败 | > 10 次/分钟 | 严重 | 运维团队 + DBA |
| 服务不可用 | 健康检查失败 | 严重 | 运维团队 + 开发团队 |

#### 告警通知方式

- **企业微信**：实时告警通知
- **邮件**：告警汇总报告
- **短信**：严重告警（凌晨或节假日）

#### 告警处理流程

1. **接收告警** → 确认告警内容和级别
2. **问题定位** → 查看日志和监控面板
3. **故障恢复** → 重启服务、切换数据库、回滚代码
4. **问题复盘** → 分析根因，优化系统

## 其他关注点

### 安全性 (Security)

#### 认证和授权

- **JWT Token 认证**：所有业务接口需要携带有效 Token
- **密码加密**：使用 bcrypt 算法，不存储明文密码
- **Token 有效期**：**强烈建议生产环境缩短为 1-24 小时，并实现 Refresh Token 机制**。当前开发环境配置的 999 小时仅用于测试便利性，生产环境绝不可使用如此长的有效期

#### 输入验证

- 参数类型验证（Gin 的 Binding 机制）
- SQL 注入防护（GORM 参数化查询）
- XSS 防护（前端负责）

#### 数据安全

- 敏感数据加密存储（用户密码）
- HTTPS 传输（生产环境必须）
- 日志脱敏（不记录密码、Token 等敏感信息）

#### 安全建议

1. **生产环境使用环境变量存储密钥**：
   - JWT 签名密钥必须使用强随机密钥（至少 256 位）
   - 通过环境变量或密钥管理服务配置，绝不可硬编码
   - 定期轮换密钥（建议每 3-6 个月）
2. **优化 Token 有效期**：
   - Access Token 有效期设置为 1-24 小时
   - 实现 Refresh Token 机制（有效期 7-30 天）
   - Refresh Token 使用单独的密钥和存储
3. **定期更新依赖包**：修复已知安全漏洞
4. **接口限流**：防止暴力破解和 DDoS 攻击
5. **审计日志**：记录关键操作，便于追溯

### 容灾和降级

#### 过载保护

- **接口限流**：使用令牌桶或漏桶算法限制 QPS
- **熔断机制**：数据库连接失败时快速失败，避免雪崩
- **超时控制**：数据库查询、HTTP 请求设置超时时间

#### 有损降级

当系统压力过大时，优先保证核心功能：

**核心功能**（必须保证）：
1. 用户登录
2. 商品查询
3. 订单创建

**非核心功能**（可降级）：
1. 购物车同步（可延迟）
2. 日志详细记录（可简化）
3. 复杂统计查询（可暂停）

#### 灾难恢复

1. **数据备份**：每日全量备份 + 增量备份
2. **异地容灾**：关键数据异地备份
3. **快速恢复**：准备灾难恢复脚本和文档

### 部署和运维

#### 部署方式

**推荐方案：Docker 容器化部署**

Dockerfile 示例：
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY config/config.yaml ./config/
EXPOSE 8080
CMD ["./server"]
```

**部署流程**：
1. 编译 Go 应用
2. 构建 Docker 镜像
3. 推送到镜像仓库
4. 在目标服务器拉取镜像并运行
5. 配置负载均衡（Nginx）

#### 运维策略

**日常运维**：
- 监控系统健康状态
- 定期检查日志和告警
- 数据库性能优化（慢查询分析）

**版本发布**：
- 灰度发布：先发布到 10% 节点观察
- 快速回滚：保留上一版本镜像
- 发布时段：业务低峰期（凌晨 2-4 点）

**应急预案**：
- 数据库故障：切换到从库
- 应用故障：重启或回滚
- 网络故障：切换备用链路

### 可扩展性设计

#### 水平扩展

- **无状态设计**：应用不保存状态，可随意增减实例
- **负载均衡**：Nginx 或云服务负载均衡
- **会话管理**：使用 JWT，无需 Session 共享

#### 未来扩展方向

1. **微服务改造**：
   - 用户服务、商品服务、订单服务独立部署
   - 使用 gRPC 或 REST 通信

2. **缓存层**：
   - Redis 缓存热点数据（商品信息、用户信息）
   - 降低数据库压力

3. **消息队列**：
   - 异步处理订单创建、库存扣减
   - 削峰填谷

4. **搜索引擎**：
   - Elasticsearch 实现商品全文搜索
   - 支持复杂查询和聚合

5. **CDN 加速**：
   - 商品图片、静态资源使用 CDN
   - 减少服务器带宽压力

### 测试策略

#### 单元测试

- **覆盖率目标**：> 70%
- **测试重点**：Service 层业务逻辑
- **Mock 工具**：`gomock` 或 `testify/mock`

#### 集成测试

- **测试数据库**：使用 SQLite 或测试数据库
- **测试场景**：完整的业务流程（注册→登录→下单）

#### 性能测试

- **工具**：`go test -bench`、`ab`、`wrk`
- **目标**：验证 SLA 指标

## 系统迁移影响

### 当前状态

本系统为新建项目，不涉及老系统迁移。

### 未来迁移考虑

如果未来需要从单体应用迁移到微服务：

1. **数据迁移**：
   - 按业务领域拆分数据库
   - 使用 ETL 工具迁移历史数据

2. **服务拆分**：
   - 先拆分读流量，再拆分写流量
   - 双写验证数据一致性

3. **灰度发布**：
   - 通过路由规则逐步切换流量
   - 保留回滚机制

4. **监控和验证**：
   - 实时监控业务指标和错误率
   - 对比新老系统的性能和准确性

## 附录

### 工程链接

- **代码仓库**：https://github.com/xingexin/server
- **项目分支**：`master`（主分支）

### 技术文档

- **Gin 文档**：https://gin-gonic.com/docs/
- **GORM 文档**：https://gorm.io/docs/
- **JWT 规范**：https://jwt.io/introduction
- **Go 最佳实践**：https://go.dev/doc/effective_go

### 相关文档

- 项目结构：`PROJECT_STRUCTURE.md`
- 项目说明：`Readme.md`
- 配置文件：`config/config.yaml`

### 更新记录

| 版本 | 日期 | 作者 | 更新内容 |
|------|------|------|---------|
| v1.0 | 2024-10-15 | Copilot | 初始版本，完成系统设计文档 |

---

**文档状态**：已完成初稿，待评审

**下一步行动**：
1. 团队 Review 设计文档
2. 确认技术选型和架构设计
3. 细化开发任务和时间计划
4. 开始编码实现
