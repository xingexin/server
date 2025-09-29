```
go-goods-management/
├── cmd/                     # 程序入口
│   └── main.go              # 启动文件
│
├── config/                  # 配置
│   └── config.yaml          # 数据库/服务配置
│   └── config.go            # 读取配置
│
├── internal/                # 内部应用逻辑
│   ├── product/             # 商品模块
│   │   ├── model.go         # 商品实体定义 (GORM 模型)
│   │   ├── repository.go    # 数据库操作 (DAO 层)
│   │   ├── service.go       # 业务逻辑 (Service 层)
│   │   └── handler.go       # 控制器 (Gin Handler)
│   │
│   ├── middleware/          # Gin 中间件 (日志/鉴权/跨域等)
│   │   └── auth.go
│   │
│   └── router/              # 路由定义
│       └── router.go
│
├── pkg/                     # 可复用工具库
│   ├── db/                  # 数据库初始化 (GORM)
│   │   └── db.go
│   └── logger/              # 日志工具
│       └── logger.go
│
├── go.mod
├── go.sum
└── README.md
```