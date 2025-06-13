package handlers

import (
	"context"
    "encoding/json"
    "net/http"
	"database/sql"
    "strconv"
	"errors"
	"time"
	"fmt"

    //"github.com/go-redis/redis/v8"
    "github.com/gorilla/mux"
    //"github.com/jmoiron/sqlx"
    //"github.com/nats-io/nats.go"
)

var _ = mux.Vars

func (h *GoodsHandler) ReprioritizeGood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var input struct {
		NewPriority int `json:"newPriority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Начинаем транзакцию
	tx, err := h.db.Beginx()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	// Получаем текущий приоритет товара
	var currentPriority, projectID int
	err = tx.Get(&currentPriority, "SELECT priority, project_id FROM goods WHERE id = $1 FOR UPDATE", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Good not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Обновляем приоритеты
	if input.NewPriority < currentPriority {
		// Увеличиваем приоритеты у записей между новым и старым приоритетом
		_, err = tx.Exec(
			"UPDATE goods SET priority = priority + 1 WHERE project_id = $1 AND priority >= $2 AND priority < $3 AND id != $4",
			projectID, input.NewPriority, currentPriority, id,
		)
	} else if input.NewPriority > currentPriority {
		// Уменьшаем приоритеты у записей между старым и новым приоритетом
		_, err = tx.Exec(
			"UPDATE goods SET priority = priority - 1 WHERE project_id = $1 AND priority > $2 AND priority <= $3 AND id != $4",
			projectID, currentPriority, input.NewPriority, id,
		)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Обновляем приоритет текущей записи
	_, err = tx.Exec("UPDATE goods SET priority = $1 WHERE id = $2", input.NewPriority, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Получаем обновленные приоритеты
	var updatedPriorities []struct {
		ID       int `db:"id"`
		Priority int `db:"priority"`
	}
	err = tx.Select(&updatedPriorities,
		"SELECT id, priority FROM goods WHERE project_id = $1 AND priority BETWEEN $2 AND $3",
		projectID, min(input.NewPriority, currentPriority), max(input.NewPriority, currentPriority),
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Завершаем транзакцию
	if err = tx.Commit(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Инвалидируем кеш
	h.invalidateCache(context.Background(), projectID)

	// Отправляем лог в ClickHouse через NATS
	h.sendLogToClickHouse("reprioritize", id, projectID, input.NewPriority)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"priorities": updatedPriorities,
	})
}

func (h *GoodsHandler) invalidateCache(ctx context.Context, projectID int) {
	keys, err := h.redis.Keys(ctx, fmt.Sprintf("goods:project:%d:*", projectID)).Result()
	if err == nil {
		for _, key := range keys {
			h.redis.Del(ctx, key)
		}
	}
}

func (h *GoodsHandler) sendLogToClickHouse(action string, goodID, projectID, priority int) {
	logData := map[string]interface{}{
		"action":     action,
		"good_id":    goodID,
		"project_id": projectID,
		"priority":   priority,
		"timestamp":  time.Now().Unix(),
	}
	
	jsonData, _ := json.Marshal(logData)
	h.nats.Publish("goods.logs", jsonData)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}