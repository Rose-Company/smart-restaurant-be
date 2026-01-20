package models

import (
	"time"
)

type Review struct {
	ID                 int                 `json:"id" gorm:"primaryKey"`
	OrderID            int                 `json:"order_id" gorm:"column:order_id"`
	Order              *Order              `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	UserID             string              `json:"user_id" gorm:"column:user_id"`
	User               *User               `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Rating             int                 `json:"rating" gorm:"column:rating"`
	Comment            string              `json:"comment" gorm:"column:comment;type:text"`
	Photos             []ReviewPhoto       `json:"photos" gorm:"foreignKey:ReviewID"`
	ItemsReviewed      []ReviewItem        `json:"items_reviewed" gorm:"foreignKey:ReviewID"`
	RestaurantResponse *RestaurantResponse `json:"restaurant_response,omitempty" gorm:"foreignKey:ReviewID"`
	Status             string              `json:"status" gorm:"column:status;default:'published'"`
	CreatedAt          time.Time           `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time           `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Review) TableName() string {
	return "customer_reviews"
}

type ReviewPhoto struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	ReviewID     int       `json:"review_id" gorm:"column:review_id"`
	URL          string    `json:"url" gorm:"column:url"`
	ThumbnailURL string    `json:"thumbnail_url" gorm:"column:thumbnail_url"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (ReviewPhoto) TableName() string {
	return "review_photos"
}

type ReviewItem struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	ReviewID    int       `json:"review_id" gorm:"column:review_id"`
	MenuItemID  int       `json:"menu_item_id" gorm:"column:menu_item_id"`
	MenuItem    *MenuItem `json:"menu_item,omitempty" gorm:"foreignKey:MenuItemID"`
	ItemRating  int       `json:"item_rating" gorm:"column:item_rating"`
	ItemComment string    `json:"item_comment" gorm:"column:item_comment;type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (ReviewItem) TableName() string {
	return "review_items"
}

type RestaurantResponse struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	ReviewID    int       `json:"review_id" gorm:"column:review_id"`
	Message     string    `json:"message" gorm:"column:message;type:text"`
	RespondedBy string    `json:"responded_by" gorm:"column:responded_by"`
	RespondedAt time.Time `json:"responded_at" gorm:"column:responded_at"`
}

func (RestaurantResponse) TableName() string {
	return "restaurant_responses"
}

// Request/Response models for reviews

type CreateReviewRequest struct {
	OrderID       int                  `json:"order_id" binding:"required"`
	Rating        int                  `json:"rating" binding:"required,min=1,max=5"`
	Comment       string               `json:"comment" binding:"required"`
	Photos        []ReviewPhotoRequest `json:"photos,omitempty"`
	ItemsReviewed []ReviewItemRequest  `json:"items_reviewed,omitempty"`
}

type ReviewItemRequest struct {
	MenuItemID  int    `json:"menu_item_id" binding:"required"`
	ItemRating  int    `json:"item_rating" binding:"min=1,max=5"`
	ItemComment string `json:"item_comment"`
}

type ReviewPhotoRequest struct {
	URL          string `json:"url" binding:"required"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type ReviewListResponse struct {
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Items    []ReviewDetailResponse `json:"items"`
	Extra    ReviewExtraInfo        `json:"extra,omitempty"`
}

type ReviewDetailResponse struct {
	ID                 int                     `json:"id"`
	OrderID            int                     `json:"order_id"`
	OrderNumber        string                  `json:"order_number"`
	Rating             int                     `json:"rating"`
	Comment            string                  `json:"comment"`
	Photos             []ReviewPhotoResponse   `json:"photos"`
	ItemsReviewed      []ReviewItemResponse    `json:"items_reviewed"`
	RestaurantResponse *RestaurantResponseView `json:"restaurant_response,omitempty"`
	Status             string                  `json:"status"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
}

type ReviewPhotoResponse struct {
	ID           int    `json:"id"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type ReviewItemResponse struct {
	MenuItemID    int    `json:"menu_item_id"`
	MenuItemName  string `json:"menu_item_name"`
	MenuItemImage string `json:"menu_item_image"`
	ItemRating    int    `json:"item_rating,omitempty"`
	ItemComment   string `json:"item_comment,omitempty"`
}

type RestaurantResponseView struct {
	Message     string    `json:"message"`
	RespondedBy string    `json:"responded_by"`
	RespondedAt time.Time `json:"responded_at"`
}

type ReviewExtraInfo struct {
	AverageRating   float64         `json:"average_rating"`
	TotalReviews    int64           `json:"total_reviews"`
	RatingBreakdown RatingBreakdown `json:"rating_breakdown,omitempty"`
}

type RatingBreakdown struct {
	FiveStar  int64 `json:"5_star"`
	FourStar  int64 `json:"4_star"`
	ThreeStar int64 `json:"3_star"`
	TwoStar   int64 `json:"2_star"`
	OneStar   int64 `json:"1_star"`
}
