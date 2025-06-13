package models

import "time"

type Good struct {
	ID          int       `json:"id" db:"id"`
	ProjectID   int       `json:"projectId" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Priority    int       `json:"priority" db:"priority"`
	Removed     bool      `json:"removed" db:"removed"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type GoodsListResponse struct {
	Meta  Meta   `json:"meta"`
	Goods []Good `json:"goods"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type ReprioritizeResponse struct {
	Priorities []Priority `json:"priorities"`
}

type Priority struct {
	ID       int `json:"id"`
	Priority int `json:"priority"`
}