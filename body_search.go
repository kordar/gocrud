package gocrud

import (
	"context"
	"errors"
	"github.com/kordar/gormext"
)

type SearchBody struct {
	Page       gormext.StrInt         `json:"page" form:"page"`
	PageSize   gormext.StrInt         `json:"pageSize" form:"pageSize"`
	Data       map[string]interface{} `json:"data,omitempty" form:"data,omitempty"`             // 数据
	Conditions []Condition            `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Sorts      []Sort                 `json:"sorts,omitempty" form:"sorts,omitempty"`           // 排序
	*CommonBody
}

func NewSearchBody(driver string, ctx context.Context) SearchBody {
	return SearchBody{
		Page:       1,
		PageSize:   15,
		Data:       make(map[string]interface{}),
		Conditions: make([]Condition, 0),
		CommonBody: NewCommonBody(driver, ctx),
	}
}

func (search *SearchBody) where(db interface{}, params map[string]string) interface{} {
	params = search.LoadDriverName(params)
	if search.Conditions != nil {
		for _, exec := range search.Conditions {
			db = exec.Where(db, params)
		}
	}
	return db
}

func (search *SearchBody) order(db interface{}, params map[string]string) interface{} {
	params = search.LoadDriverName(params)
	if search.Sorts != nil {
		for _, exec := range search.Sorts {
			db = exec.Order(db, params)
		}
	}
	return db
}

func (search *SearchBody) Query(db interface{}, params map[string]string) interface{} {
	db = search.where(db, params)
	db = search.order(db, params)
	return db
}

func (search *SearchBody) Paginate(db interface{}, params map[string]string) (interface{}, error) {
	db = search.Query(db, params)
	offset := (search.Page - 1) * search.PageSize
	exec := GetExecute("PAGE", search.DriverName(params), "")
	if exec == nil {
		return db, errors.New("execution function for 'PAGE' not found")
	}
	return exec(db, "", offset, search.PageSize), nil
}

func (search *SearchBody) QueryCustom(f func(search *SearchBody) interface{}) interface{} {
	return f(search)
}
