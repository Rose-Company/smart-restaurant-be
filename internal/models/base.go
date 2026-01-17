package models

import (
	"app-noti/common"
	"strings"
)

type QuerySort struct {
	Origin string
}

// Parse the query string to order string (Ex: http://example.com/messages?sort=created_at.asc,updated_at.acs
// => order string: created_at asc,updated_at acs created_at desc)
func (s QuerySort) Parse() string {
	return strings.ReplaceAll(s.Origin, ".", " ")
}

type QueryParams struct {
	Limit  int
	Offset int
	QuerySort
	Preload  []common.Preload
	Selected []string
}

type APIResponse[T any] struct {
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message"`
}

type VNPayResponse struct {
	Message    string `json:"message"`
	PaymentURL string `json:"paymentUrl,omitempty"`
	Code       string `json:"code"`
	Amount     int64  `json:"amount"`
}

type BaseRequestParamsUri struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Sort     string `form:"sort"`
}

type BaseListResponse struct {
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Items    interface{} `json:"items"`
	Extra    interface{} `json:"extra"`
}
