# 项目结构

- 仓库：server
- 生成时间：2025-09-29 14:01:22 UTC
- 深度：3
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
│   │   ├── repository/
│   │   └── service/
│   └── router/
│       └── router.go
├── main.go
└── pkg/
    ├── db/
    │   └── dataBase.go
    └── logger/
        └── logger.go

14 directories, 13 files
```

> 本文件由 GitHub Actions 自动生成，请勿手动编辑。
