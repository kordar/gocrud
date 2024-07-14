package gocrud

import (
	"encoding/json"
	"errors"
)

type commonBody struct {
	Driver string
}

func (common *commonBody) GetDriver(parallel map[string]string) string {
	if parallel == nil || parallel["driver"] == "" {
		return common.Driver
	}
	return parallel["driver"]
}

type FormBody struct {
	Conditions  []condition `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Object      interface{} `json:"object,omitempty" form:"object,omitempty"`
	safeCounter int         // 防止无条件更新
	commonBody
}

func NewFormBody(driver string) FormBody {
	return FormBody{
		commonBody: commonBody{Driver: driver},
		Conditions: make([]condition, 0),
	}
}

func (form *FormBody) GetObject(target interface{}) error {
	if marshal, err := json.Marshal(form.Object); err != nil {
		return err
	} else {
		return json.Unmarshal(marshal, target)
	}
}

func (form *FormBody) where(db interface{}, parallel map[string]string) interface{} {
	for _, exec := range form.Conditions {
		db = exec.Where(db, parallel)
	}
	return db
}

func (form *FormBody) Query(db interface{}, parallel map[string]string) interface{} {
	return form.where(db, parallel)
}

func (form *FormBody) whereSafe(db interface{}, parallel map[string]string) interface{} {
	for _, exec := range form.Conditions {
		flag := false
		db, flag = exec.WhereSafe(db, parallel)
		if flag == true {
			form.safeCounter++
		}
	}
	return db
}

func (form *FormBody) QuerySafe(db interface{}, parallel map[string]string) (interface{}, error) {
	db = form.whereSafe(db, parallel)
	if form.safeCounter == 0 {
		return db, errors.New("forbid no condition update")
	}
	return db, nil
}

func (form *FormBody) QueryCustom(
	db interface{},
	parallel map[string]string,
	fun func(form *FormBody, db interface{}, parallel map[string]string) interface{},
) interface{} {
	return fun(form, db, parallel)
}

// Create 创建model
func (form *FormBody) Create(model interface{}, db interface{}, parallel map[string]string) (interface{}, error) {
	err := form.GetObject(model)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("CREATE", form.GetDriver(parallel), "")
	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}

// CreateWithValid 创建并且校验提交参数
func (form *FormBody) CreateWithValid(model interface{}, db interface{}, parallel map[string]string, valid func(model interface{}) error) (interface{}, error) {
	err := form.GetObject(model)
	if err != nil {
		return nil, err
	}

	err = valid(model)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("CREATE", form.GetDriver(parallel), "")
	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}

// Update 更新
func (form *FormBody) Update(model interface{}, db interface{}, parallel map[string]string) (interface{}, error) {

	err := form.GetObject(model)
	if err != nil {
		return nil, err
	}

	db, err = form.QuerySafe(db, nil)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("UPDATES", form.GetDriver(parallel), "")
	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}

// UpdateWithValid 更新并且校验提交参数
func (form *FormBody) UpdateWithValid(model interface{}, db interface{}, parallel map[string]string, valid func(model interface{}) error) (interface{}, error) {

	err := form.GetObject(model)
	if err != nil {
		return nil, err
	}

	err = valid(model)
	if err != nil {
		return nil, err
	}

	db, err = form.QuerySafe(db, nil)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("UPDATES", form.GetDriver(parallel), "")
	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}

// UpdateMapWithValid 更新并且校验提交参数
func (form *FormBody) UpdateMapWithValid(model interface{}, db interface{}, parallel map[string]string, valid func(model interface{}) (error, map[string]interface{})) (interface{}, error) {

	err := form.GetObject(model)
	if err != nil {
		return nil, err
	}

	err, m := valid(model)
	if err != nil {
		return nil, err
	}

	db, err = form.QuerySafe(db, nil)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("UPDATES", form.GetDriver(parallel), "")
	if e := exec(db, "", model, m); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}

// SaveWithValid 更新并且校验提交参数
func (form *FormBody) SaveWithValid(model interface{}, db interface{}, parallel map[string]string, valid func(model interface{}) error) (interface{}, error) {

	err := form.GetObject(model)
	if err != nil {
		return nil, err
	}

	err = valid(model)
	if err != nil {
		return nil, err
	}

	db, err = form.QuerySafe(db, nil)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("SAVE", form.GetDriver(parallel), "")
	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}
