# 项目结构

- 仓库：server
- 生成时间：2025-12-11 05:16:09 UTC
- 深度：99
- 忽略：.git|target|node_modules|.idea|.vscode|dist|build

```text

├── .github/
│   └── workflows/
│       └── generate-structure.yml
├── .gitignore
├── DESIGN.md
├── PROJECT_STRUCTURE.md
├── Readme.md
├── config/
│   └── config.go
├── go.mod
├── go.sum
├── internal/
│   ├── middleware/
│   │   └── auth.go
│   ├── product/
│   │   ├── cart/
│   │   │   ├── dto/
│   │   │   │   ├── request.go
│   │   │   │   └── response.go
│   │   │   ├── handler/
│   │   │   │   └── cart_handler.go
│   │   │   ├── model/
│   │   │   │   └── cart_model.go
│   │   │   ├── repository/
│   │   │   │   └── cart_repository.go
│   │   │   └── service/
│   │   │       └── cart_service.go
│   │   ├── commodity/
│   │   │   ├── dto/
│   │   │   │   ├── request.go
│   │   │   │   └── response.go
│   │   │   ├── handler/
│   │   │   │   └── commodity_handler.go
│   │   │   ├── model/
│   │   │   │   └── commodity_model.go
│   │   │   ├── repository/
│   │   │   │   ├── commodity_repository.go
│   │   │   │   └── stock_cache_repository.go
│   │   │   └── service/
│   │   │       ├── commodity_service.go
│   │   │       └── stock_cache_service.go
│   │   ├── order/
│   │   │   ├── dto/
│   │   │   │   ├── request.go
│   │   │   │   └── response.go
│   │   │   ├── handler/
│   │   │   │   └── order_handler.go
│   │   │   ├── model/
│   │   │   │   └── order_model.go
│   │   │   ├── repository/
│   │   │   │   ├── my_errors.go
│   │   │   │   ├── order_dq_repository.go
│   │   │   │   └── order_repository.go
│   │   │   └── service/
│   │   │       ├── order_cancel_service.go
│   │   │       └── order_service.go
│   │   ├── scheduler/
│   │   │   ├── order_dq_scheduler.go
│   │   │   ├── recovery_scheduler.go
│   │   │   └── scheduler.go
│   │   └── user/
│   │       ├── dto/
│   │       │   ├── request.go
│   │       │   └── response.go
│   │       ├── handler/
│   │       │   └── user_handler.go
│   │       ├── model/
│   │       │   └── user_model.go
│   │       ├── repository/
│   │       │   └── user_repository.go
│   │       └── service/
│   │           └── user_service.go
│   └── router/
│       └── router.go
├── main.go
└── pkg/
    ├── container/
    │   └── container.go
    ├── db/
    │   └── dataBase.go
    ├── logger/
    │   └── logger.go
    ├── redis/
    │   └── redis.go
    └── response/
        ├── code.go
        └── response.go

39 directories, 49 files
```

> 本文件由 GitHub Actions 自动生成，请勿手动编辑。
