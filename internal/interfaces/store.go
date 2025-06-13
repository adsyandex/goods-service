package interfaces

import (
	"context"
	//"database/sql"
	"time"
	
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	//"github.com/nats-io/nats.go"
)

// Database интерфейс
type Database interface {
	sqlx.Ext
	Beginx() (*sqlx.Tx, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// Redis интерфейс
type Redis interface {
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

// NATS интерфейс
type NATS interface {
	Publish(subject string, data []byte) error
}