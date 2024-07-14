package gocrud

type SearchBody struct {
	Page       int                    `json:"page" form:"page"`
	PageSize   int                    `json:"pageSize" form:"pageSize"`
	Data       map[string]interface{} `json:"data,omitempty" form:"data,omitempty"`             // 数据
	Conditions []condition            `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Sorts      []sort                 `json:"sorts,omitempty" form:"sorts,omitempty"`           // 排序
	commonBody
}

func NewSearchBody(driver string) SearchBody {
	return SearchBody{
		Page:       1,
		PageSize:   15,
		Data:       make(map[string]interface{}),
		Conditions: make([]condition, 0),
		commonBody: commonBody{Driver: driver},
	}
}

func (search *SearchBody) where(db interface{}, parallel map[string]string) interface{} {
	for _, exec := range search.Conditions {
		db = exec.Where(db, parallel)
	}
	return db
}

func (search *SearchBody) order(db interface{}, parallel map[string]string) interface{} {
	for _, exec := range search.Sorts {
		db = exec.Order(db, parallel)
	}
	return db
}

func (search *SearchBody) Paginate(db interface{}, parallel map[string]string) interface{} {
	db = search.Query(db, parallel)
	offset := (search.Page - 1) * search.PageSize
	execute := GetExecute("PAGE", search.GetDriver(parallel), "")
	return execute(db, "", offset, search.PageSize)
}

func (search *SearchBody) Query(db interface{}, parallel map[string]string) interface{} {
	db = search.where(db, parallel)
	db = search.order(db, parallel)
	return db
}

func (search *SearchBody) QueryCustom(f func(search *SearchBody) interface{}) interface{} {
	return f(search)
}
