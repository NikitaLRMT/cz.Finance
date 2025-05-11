package models

import (
	"time"
)

// WishlistPriority перечисляет возможные приоритеты для элементов списка желаний
type WishlistPriority string

const (
	PriorityHigh   WishlistPriority = "high"
	PriorityMedium WishlistPriority = "medium"
	PriorityLow    WishlistPriority = "low"
)

// WishlistItem представляет модель элемента списка желаний
type WishlistItem struct {
	ID          int64            `json:"id" db:"id"`
	UserID      int64            `json:"user_id" db:"user_id"`
	Title       string           `json:"title" db:"title" validate:"required,min=2,max=100"`
	Price       float64          `json:"price" db:"price" validate:"required,gt=0"`
	Priority    WishlistPriority `json:"priority" db:"priority" validate:"required"`
	Description string           `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

// CreateWishlistItemRequest модель для создания нового элемента списка желаний
type CreateWishlistItemRequest struct {
	Title       string           `json:"title" validate:"required,min=2,max=100"`
	Price       float64          `json:"price" validate:"required,gt=0"`
	Priority    WishlistPriority `json:"priority" validate:"required"`
	Description string           `json:"description,omitempty"`
}

// UpdateWishlistItemRequest модель для обновления элемента списка желаний
type UpdateWishlistItemRequest struct {
	Title       *string           `json:"title" validate:"omitempty,min=2,max=100"`
	Price       *float64          `json:"price" validate:"omitempty,gt=0"`
	Priority    *WishlistPriority `json:"priority"`
	Description *string           `json:"description,omitempty"`
}
