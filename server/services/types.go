package services

import "time"

type Group struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
	IsPublic    bool      `json:"is_public"`
}

// resembles todo_app.items
type Item struct {
	Id        int        `json:"id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	RemindAt  *time.Time `json:"remind_at"`
	IsActive  bool       `json:"is_active"`
}

type GroupedItem struct {
	Id        int       `json:"id"`
	GroupId   int       `json:"group_id"`
	ItemId    int       `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

type Content struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Updatedat time.Time `json:"updated_at"`
}

type ItemContentResponse struct {
	ItemId  int    `json:"item_id"`
	Id      int    `json:"id"`
	Content string `json:"content"`
}
