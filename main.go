package main

import (
	"log"
	"net/http"
	"os"
	
	"github.com/joho/godotenv"
	"goods-service/internal/app"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func main() {
	  // Загрузка .env
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
	
	// Инициализация подключений
	 dbDSN := os.Getenv("DB_DSN")
    db, err := sqlx.Connect("postgres", dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	natsConn, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	// Создание приложения
	application := app.NewApp(db, redisClient, natsConn)

	// Настройка маршрутов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Goods Service is running"))
	})
	application.SetupRoutes()

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}