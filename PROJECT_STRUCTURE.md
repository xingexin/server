# 项目结构

- 仓库：server
- 生成时间：2025-10-11 05:45:09 UTC
- 深度：99
- 忽略：.git|target|node_modules|.idea|.vscode|dist|build

```text

├── .github/
│   └── workflows/
│       └── generate-structure.yml
├── .gitignore
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

14 directories, 20 files
```

> 本文件由 GitHub Actions 自动生成，请勿手动编辑。
