# goods-service

## Структура проекта
```
/goods-service
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── migrations
│   ├── 001_init_tables.up.sql
│   └── 002_clickhouse_logs.up.sql
└── internal
    ├── app
    │   ├── app.go
    │   └── server.go
    ├── models
    │   ├── good.go
    │   └── project.go
    └── handlers
        ├── goods.go
        └── reprioritize.go
```
