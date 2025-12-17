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
