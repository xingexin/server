# 项目结构

- 仓库：server
- 生成时间：2025-10-20 05:41:31 UTC
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
│   ├── config.go
│   └── config.yaml
├── go.mod
├── go.sum
├── internal/
│   ├── middleware/
│   │   └── auth.go
│   ├── product/
│   │   ├── handler/
│   │   │   ├── cart_handler.go
│   │   │   ├── commodity_handler.go
│   │   │   ├── order_handler.go
│   │   │   └── user_handler.go
│   │   ├── model/
│   │   │   ├── cart_model.go
│   │   │   ├── commodity_model.go
│   │   │   ├── order_model.go
│   │   │   └── user_model.go
│   │   ├── repository/
│   │   │   ├── cart_repository.go
│   │   │   ├── commodity_repository.go
│   │   │   ├── order_repository.go
│   │   │   └── user_repository.go
│   │   └── service/
│   │       ├── cart_service.go
│   │       ├── commodity_service.go
│   │       ├── order_service.go
│   │       └── user_service.go
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
    └── response/
        ├── code.go
        └── response.go

17 directories, 33 files
```

> 本文件由 GitHub Actions 自动生成，请勿手动编辑。
