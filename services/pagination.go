package services

import (
	"math"

	"gorm.io/gorm"
)

// Pagination 分页参数
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	Pages    int   `json:"pages"`
}

// NewPagination 创建分页参数（自动修正边界）
func NewPagination(page, pageSize int) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return &Pagination{Page: page, PageSize: pageSize}
}

// Offset 计算偏移量
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// SetTotal 设置总数并计算总页数
func (p *Pagination) SetTotal(total int64) {
	p.Total = total
	p.Pages = int(math.Ceil(float64(total) / float64(p.PageSize)))
}

// Paginate GORM scope，用于链式调用
//
// 用法：
//
//	p := services.NewPagination(page, pageSize)
//	db.Model(&Model{}).Scopes(p.Paginate()).Find(&results)
func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.Offset()).Limit(p.PageSize)
	}
}

// PaginateQuery 分页查询辅助（一步完成 count + find）
//
// 用法：
//
//	var users []models.User
//	p, err := services.PaginateQuery(db.Model(&models.User{}).Where("status = ?", "active"), page, pageSize, &users)
func PaginateQuery(query *gorm.DB, page, pageSize int, dest interface{}) (*Pagination, error) {
	p := NewPagination(page, pageSize)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	p.SetTotal(total)

	if err := query.Offset(p.Offset()).Limit(p.PageSize).Find(dest).Error; err != nil {
		return nil, err
	}
	return p, nil
}
