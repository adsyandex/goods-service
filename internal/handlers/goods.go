package handlers

import (
	//"database/sql"
	"encoding/json"
	"net/http"
	
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
)

type GoodsHandler struct {
	db    *sqlx.DB
	redis *redis.Client
	nats  *nats.Conn
}

func NewGoodsHandler(db *sqlx.DB, redis *redis.Client, nats *nats.Conn) *GoodsHandler {
	return &GoodsHandler{
		db:    db,
		redis: redis,
		nats:  nats,
	}
}

// CreateGood - обработчик создания товара
func (h *GoodsHandler) CreateGood(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Реализация создания товара в БД
	_, err := h.db.Exec("INSERT INTO goods (name, description) VALUES ($1, $2)", 
		input.Name, input.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateGood - обработчик обновления товара
func (h *GoodsHandler) UpdateGood(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec(
		"UPDATE goods SET name = $1, description = $2 WHERE id = $3",
		input.Name, input.Description, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteGood - обработчик удаления товара
func (h *GoodsHandler) DeleteGood(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, err := h.db.Exec("DELETE FROM goods WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ListGoods - обработчик получения списка товаров
func (h *GoodsHandler) ListGoods(w http.ResponseWriter, r *http.Request) {
	var goods []struct {
		ID          int    `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
	}

	err := h.db.Select(&goods, "SELECT id, name, description FROM goods")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(goods)
}