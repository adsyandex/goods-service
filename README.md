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
# Goods Service 🛍️

Микросервис для управления товарами с поддержкой CRUD операций, реализованный на Go.

## 🚀 Особенности

- Полный CRUD функционал (создание, чтение, обновление, удаление товаров)
- Поддержка приоритизации товаров
- Кеширование через Redis
- Логирование операций в ClickHouse через NATS
- RESTful API
- Контейнеризация с Docker

## 📦 Технологический стек

- **Язык**: Go 1.24.1
- **База данных**: PostgreSQL
- **Кеширование**: Redis
- **Очереди**: NATS
- **Хранилище логов**: ClickHouse
- **Роутинг**: gorilla/mux

## 🔧 Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/goods-service.git
cd goods-service
```
## Запуск через Docker
```
docker-compose up --build
```

## Пример запроса
#### Создание запроса:
```
curl -X POST http://localhost:8080/goods \
  -H "Content-Type: application/json" \
  -d '{"name": "Новый товар", "description": "Описание товара"}'
```


