package gocrud

import "errors"

type SearchBody struct {
	Page       int                    `json:"page" form:"page"`
	PageSize   int                    `json:"pageSize" form:"pageSize"`
	Data       map[string]interface{} `json:"data,omitempty" form:"data,omitempty"`             // 数据
	Conditions []Condition            `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Sorts      []Sort                 `json:"sorts,omitempty" form:"sorts,omitempty"`           // 排序
	*CommonBody
}

func NewSearchBody(driver string, ctx interface{}) SearchBody {
	return SearchBody{
		Page:       1,
		PageSize:   15,
		Data:       make(map[string]interface{}),
		Conditions: make([]Condition, 0),
		CommonBody: NewCommonBody(driver, ctx),
	}
}

func (search *SearchBody) where(db interface{}, parallel map[string]string) interface{} {
	parallel = search.LoadDriverName(parallel)
	if search.Conditions != nil {
		for _, exec := range search.Conditions {
			db = exec.Where(db, parallel)
		}
	}
	return db
}

func (search *SearchBody) order(db interface{}, parallel map[string]string) interface{} {
	parallel = search.LoadDriverName(parallel)
	if search.Sorts != nil {
		for _, exec := range search.Sorts {
			db = exec.Order(db, parallel)
		}
	}
	return db
}

func (search *SearchBody) Query(db interface{}, parallel map[string]string) interface{} {
	db = search.where(db, parallel)
	db = search.order(db, parallel)
	return db
}

func (search *SearchBody) Paginate(db interface{}, parallel map[string]string) (interface{}, error) {
	db = search.Query(db, parallel)
	offset := (search.Page - 1) * search.PageSize
	exec := GetExecute("PAGE", search.DriverName(parallel), "")
	if exec == nil {
		return db, errors.New("execution function for 'PAGE' not found")
	}
	return exec(db, "", offset, search.PageSize), nil
}

func (search *SearchBody) QueryCustom(f func(search *SearchBody) interface{}) interface{} {
	return f(search)
}
