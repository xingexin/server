# 项目结构

- 仓库：server
- 生成时间：2025-10-09 12:41:18 UTC
- 深度：99
- 忽略：.git|target|node_modules|.idea|.vscode|dist|build

```text

├── .github/
│   └── workflows/
│       └── generate-structure.yml
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
│   │   ├── handler.go
│   │   ├── model/
│   │   │   ├── commodity_model.go
│   │   │   └── user_model.go
│   │   ├── repository/
│   │   │   ├── commodity_repository.go
│   │   │   └── user_repository.go
│   │   └── service/
│   │       ├── commodity_service.go
│   │       └── user_service.go
│   └── router/
│       └── router.go
├── main.go
└── pkg/
    ├── db/
    │   └── dataBase.go
    └── logger/
        └── logger.go

14 directories, 19 files
```

> 本文件由 GitHub Actions 自动生成，请勿手动编辑。
