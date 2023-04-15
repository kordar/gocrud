package gocrud

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SearchBody struct {
	Page       int                    `json:"page" form:"page"`
	PageSize   int                    `json:"pageSize" form:"pageSize"`
	Data       map[string]interface{} `json:"data,omitempty" form:"data,omitempty"`             // 数据
	Conditions []*condition           `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Sorts      []*sort                `json:"sorts,omitempty" form:"sorts,omitempty"`           // 排序
	commonBody
}

func NewSearchBody(ctx *gin.Context) SearchBody {
	return SearchBody{
		Page:       1,
		PageSize:   15,
		Data:       make(map[string]interface{}),
		Conditions: make([]*condition, 0),
		commonBody: commonBody{Ctx: ctx},
	}
}

func (search *SearchBody) where(db *gorm.DB, parallel map[string]string) *gorm.DB {
	for _, condition := range search.Conditions {
		db = condition.Where(db, parallel)
	}
	return db
}

func (search *SearchBody) order(db *gorm.DB, parallel map[string]string) *gorm.DB {
	for _, sort := range search.Sorts {
		db = sort.Order(db, parallel)
	}
	return db
}

func (search *SearchBody) Paginate(db *gorm.DB, parallel map[string]string) *gorm.DB {
	db = search.Query(db, parallel)
	offset := (search.Page - 1) * search.PageSize
	db.Offset(offset).Limit(search.PageSize)
	return db
}

func (search *SearchBody) Query(db *gorm.DB, parallel map[string]string) *gorm.DB {
	db = search.where(db, parallel)
	db = search.order(db, parallel)
	return db
}

func (search *SearchBody) QueryCustom(db *gorm.DB, parallel map[string]string, fun func(search *SearchBody, db *gorm.DB, parallel map[string]string) *gorm.DB) *gorm.DB {
	return fun(search, db, parallel)
}
