package gocrud

import "errors"

type SearchBody[T interface{}, C interface{}] struct {
	Page       int                    `json:"page" form:"page"`
	PageSize   int                    `json:"pageSize" form:"pageSize"`
	Data       map[string]interface{} `json:"data,omitempty" form:"data,omitempty"`             // 数据
	Conditions []Condition            `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Sorts      []Sort                 `json:"sorts,omitempty" form:"sorts,omitempty"`           // 排序
	*CommonBody[C]
}

func NewSearchBody[T interface{}, C interface{}](driver string, ctx C) SearchBody[T, C] {
	return SearchBody[T, C]{
		Page:       1,
		PageSize:   15,
		Data:       make(map[string]interface{}),
		Conditions: make([]Condition, 0),
		CommonBody: NewCommonBody[C](driver, ctx),
	}
}

func (search *SearchBody[T, C]) where(db T, parallel map[string]string) T {
	parallel = search.LoadDriverName(parallel)
	if search.Conditions != nil {
		for _, exec := range search.Conditions {
			db = exec.Where(db, parallel).(T)
		}
	}
	return db
}

func (search *SearchBody[T, C]) order(db T, parallel map[string]string) T {
	parallel = search.LoadDriverName(parallel)
	if search.Sorts != nil {
		for _, exec := range search.Sorts {
			db = exec.Order(db, parallel).(T)
		}
	}
	return db
}

func (search *SearchBody[T, C]) Query(db T, parallel map[string]string) T {
	db = search.where(db, parallel)
	db = search.order(db, parallel)
	return db
}

func (search *SearchBody[T, C]) Paginate(db T, parallel map[string]string) (T, error) {
	db = search.Query(db, parallel)
	offset := (search.Page - 1) * search.PageSize
	exec := GetExecute("PAGE", search.DriverName(parallel), "")
	if exec == nil {
		return db, errors.New("execution function for 'PAGE' not found")
	}
	return exec(db, "", offset, search.PageSize).(T), nil
}

func (search *SearchBody[T, C]) QueryCustom(f func(search *SearchBody[T, C]) T) T {
	return f(search)
}
