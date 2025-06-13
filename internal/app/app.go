package app

import (
	"github.com/jmoiron/sqlx"
	nats "github.com/nats-io/nats.go"
	redis "github.com/go-redis/redis/v8"
)

type App struct {
	DB     *sqlx.DB
	Redis  *redis.Client
	Nats   *nats.Conn
}

func NewApp(db *sqlx.DB, redis *redis.Client, nats *nats.Conn) *App {
	return &App{
		DB:     db,
		Redis:  redis,
		Nats:   nats,
	}
}