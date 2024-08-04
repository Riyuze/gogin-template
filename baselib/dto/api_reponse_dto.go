package dto

import (
	"math"
	"strings"
)

type Response[T any] struct {
	ResponseCode    string    `json:"responseCode" example:"200"`                           // Http Response Code
	ResponseMessage string    `json:"responseMessage" example:"Messages Messages Messages"` // Response Message
	Data            T         `json:"data"`                                                 // Data (Any model)
	LogReff         string    `json:"logReff" example:"LogReffLogReffLogReffLogReff"`       // LogReff (Use this to search in splunk)
	TraceId         string    `json:"traceId" example:"TraceIdTraceIdTraceIdTraceId"`       // TraceId (Use this as trace id in jaeger)
	PageInfo        *PageInfo `json:"pageInfo"`                                             // PageInfo (Only for response type list with pages)
}

type PageRequest struct {
	Query         string `json:"search"`        // Search Query
	Page          int    `json:"page"`          // Current Page of the Data
	PageSize      int    `json:"pageSize"`      // Amount of Items Shown on the Page
	SortBy        string `json:"sortBy"`        // Data to Sort By
	SortDirection string `json:"sortDirection"` // Sort Direction of the Data
}

type PageInfo struct {
	CurrentPageIndex    int `json:"currentPageIndex" example:"1"`       // Current Page Index
	MaxPageIndex        int `json:"maxPageIndex" example:"10"`          // Max Page Index
	RowsPerPage         int `json:"rowsPerPage" example:"100"`          // Rows Per Page
	TotalAvailableItems int `json:"totalAvailableItems" example:"1000"` // Total Available Items
}

func (p *PageRequest) GetOrderString(baseKey string, allowedOrder string) (orderString string) {
	allowedSort := strings.Split(allowedOrder, ",")
	for i := range allowedSort {
		allowedSort[i] = strings.Trim(allowedSort[i], " ")
	}
	orderString = ""
	if p.Page < 0 {
		p.Page = 0
	}
	if p.SortDirection != "ASC" && p.SortDirection != "DESC" {
		p.SortDirection = "ASC"
	}
	if p.SortBy == "" || !strings.Contains(strings.Join(allowedSort, ","), strings.ToLower(p.SortBy)) {
		orderString = baseKey
	} else {
		orderString = p.SortBy + " " + p.SortDirection + ", " + baseKey
	}
	return orderString
}

func (p *PageRequest) GetPageInfo(totalData int) (pageInfo PageInfo) {
	rowsPerPage := totalData
	if rowsPerPage <= 0 {
		rowsPerPage = 10
	}
	maxPage := int(math.Ceil(float64(totalData)/float64(rowsPerPage))) - 1
	if maxPage < 0 {
		maxPage = 0
	}
	pageInfo = PageInfo{
		CurrentPageIndex:    p.Page,
		MaxPageIndex:        maxPage,
		RowsPerPage:         rowsPerPage,
		TotalAvailableItems: totalData,
	}
	return pageInfo
}
