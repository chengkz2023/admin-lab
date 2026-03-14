# Server

这是当前脚手架的后端部分，基于 `gin + gorm + casbin`。

## 当前定位

- 只保留超级管理员相关的系统能力
- 默认数据库为 MySQL
- 默认缓存为 Redis
- 默认文件存储为本地存储
- 启动时自动建表并补齐系统出厂数据

## 目录说明

```text
server
├─ api
├─ config
├─ core
├─ global
├─ initialize
├─ middleware
├─ model
├─ resource
├─ router
├─ service
├─ source
└─ utils
```

## 运行

```bash
go run .
```

## 说明

- 当前仓库已经移除了 plugin、AI、自动代码生成运行链、多对象存储、多数据库初始化等能力
- 如果系统表为空，服务启动后会自动补系统默认数据
- Swagger 页面已经移除，如需重新接入，需要单独恢复文档生成链
