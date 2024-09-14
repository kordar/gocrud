package gocrud

import (
	"context"
	"encoding/json"
	"errors"
)

type FormBody struct {
	Conditions  []Condition `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Object      interface{} `json:"object,omitempty" form:"object,omitempty"`
	safeCounter int         // 防止无条件更新
	*CommonBody
}

func NewFormBody(driver string, ctx context.Context) FormBody {
	return FormBody{
		CommonBody: NewCommonBody(driver, ctx),
		Conditions: make([]Condition, 0),
	}
}

func (form *FormBody) Unmarshal(target interface{}) error {
	if marshal, err := json.Marshal(form.Object); err != nil {
		return err
	} else {
		return json.Unmarshal(marshal, target)
	}
}

func (form *FormBody) where(db interface{}, params map[string]string) interface{} {
	params = form.LoadDriverName(params)
	if form.Conditions != nil {
		for _, exec := range form.Conditions {
			db = exec.Where(db, params)
		}
	}
	return db
}

func (form *FormBody) whereSafe(db interface{}, params map[string]string) interface{} {
	params = form.LoadDriverName(params)
	if form.Conditions != nil {
		for _, exec := range form.Conditions {
			db2, flag := exec.WhereSafe(db, params)
			db = db2
			if flag == true {
				form.safeCounter++
			}
		}
	}
	return db
}

func (form *FormBody) Query(db interface{}, params map[string]string) interface{} {
	return form.where(db, params)
}

func (form *FormBody) QuerySafe(db interface{}, params map[string]string) (interface{}, error) {
	db = form.whereSafe(db, params)
	if form.safeCounter == 0 {
		return db, errors.New("forbid no condition update")
	}
	return db, nil
}

func (form *FormBody) QueryCustom(f func(form *FormBody) interface{}) interface{} {
	return f(form)
}

// Create 创建model
func (form *FormBody) Create(model interface{}, db interface{}, params map[string]string) (interface{}, error) {
	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("CREATE", form.DriverName(params), "")
	if exec == nil {
		return nil, errors.New("execution function for 'CREATE' not found")
	}

	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}

// CreateWithValid 创建并且校验提交参数
func (form *FormBody) CreateWithValid(model interface{}, db interface{}, params map[string]string, valid func(model interface{}) error) (interface{}, error) {
	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	err = valid(model)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("CREATE", form.DriverName(params), "")
	if exec == nil {
		return nil, errors.New("execution function for 'CREATE' not found")
	}

	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}

// Update 更新模型
func (form *FormBody) Update(model interface{}, db interface{}, params map[string]string) (interface{}, error) {

	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	// TODO 更新model必须提供有效的更新条件
	db, err = form.QuerySafe(db, params)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("UPDATES", form.DriverName(params), "")
	if exec == nil {
		return nil, errors.New("execution function for 'UPDATES' not found")
	}

	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}

// UpdateWithValid 更新并且校验提交参数
func (form *FormBody) UpdateWithValid(model interface{}, db interface{}, params map[string]string, valid func(model interface{}) error) (interface{}, error) {

	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	err = valid(model)
	if err != nil {
		return nil, err
	}

	db, err = form.QuerySafe(db, params)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("UPDATES", form.DriverName(params), "")
	if exec == nil {
		return nil, errors.New("execution function for 'UPDATES' not found")
	}

	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}

// UpdateMapWithValid 更新并且校验提交参数
func (form *FormBody) UpdateMapWithValid(model interface{}, db interface{}, params map[string]string, valid func(model interface{}) (error, map[string]interface{})) (interface{}, error) {

	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	err, m := valid(model)
	if err != nil {
		return nil, err
	}

	db, err = form.QuerySafe(db, params)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("UPDATES", form.DriverName(params), "")
	if exec == nil {
		return nil, errors.New("execution function for 'UPDATES' not found")
	}

	if e := exec(db, "", model, m); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}
}

// Save 更新并且校验提交参数
func (form *FormBody) Save(model interface{}, db interface{}, params map[string]string) (interface{}, error) {

	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("SAVE", form.DriverName(params), "")
	if exec == nil {
		return nil, errors.New("execution function for 'SAVE' not found")
	}

	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}

}

// SaveWithValid 更新并且校验提交参数
func (form *FormBody) SaveWithValid(model interface{}, db interface{}, params map[string]string, valid func(model interface{}) error) (interface{}, error) {

	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	err = valid(model)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("SAVE", form.DriverName(params), "")
	if exec == nil {
		return nil, errors.New("execution function for 'SAVE' not found")
	}

	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}

}
