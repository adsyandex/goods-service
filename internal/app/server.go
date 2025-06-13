package app

import (
	"net/http"
	"goods-service/internal/handlers"
)

func (a *App) SetupRoutes() {
	handler := handlers.NewGoodsHandler(a.DB, a.Redis, a.Nats)

	http.HandleFunc("/goods", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateGood(w, r)
		case http.MethodGet:
			handler.ListGoods(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/goods/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			handler.UpdateGood(w, r)
		case http.MethodDelete:
			handler.DeleteGood(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}