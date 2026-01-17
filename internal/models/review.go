package models

import "time"

type ReviewItem struct {
	ID          int        `json:"id"`
	ReviewId    int        `json:"review_id"`
	MenuItemId  int        `json:"menu_item_id"`
	ItemRating  int        `json:"item_rating"`
	ItemComment string     `json:"item_comment"`
	CreatedAt   *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
}
