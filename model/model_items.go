package model

import "time"

type Item struct {
	ID         int64      `json:"id"`
	CategoryID int64      `json:"category_id"`
	ItemName   string     `json:"item_name"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
