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


**依赖隔离**：
- 通过接口抽象外部依赖
- 便于测试时 Mock 外部依赖
- 降低对具体实现的耦合

## 监控和告警


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

## 附录

### 工程链接

- **代码仓库**：https://github.com/xingexin/server
- **项目分支**：`master`（主分支）
