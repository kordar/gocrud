package gocrud

import (
	"encoding/json"
	"errors"
)

type FormBody struct {
	Conditions  []Condition `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Object      interface{} `json:"object,omitempty" form:"object,omitempty"`
	safeCounter int         // 防止无条件更新
	*CommonBody
}

func NewFormBody(driver string, ctx interface{}) FormBody {
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

func (form *FormBody) where(db interface{}, parallel map[string]string) interface{} {
	parallel = form.LoadDriverName(parallel)
	if form.Conditions != nil {
		for _, exec := range form.Conditions {
			db = exec.Where(db, parallel)
		}
	}
	return db
}

func (form *FormBody) whereSafe(db interface{}, parallel map[string]string) interface{} {
	parallel = form.LoadDriverName(parallel)
	if form.Conditions != nil {
		for _, exec := range form.Conditions {
			db2, flag := exec.WhereSafe(db, parallel)
			db = db2
			if flag == true {
				form.safeCounter++
			}
		}
	}
	return db
}

func (form *FormBody) Query(db interface{}, parallel map[string]string) interface{} {
	return form.where(db, parallel)
}

func (form *FormBody) QuerySafe(db interface{}, parallel map[string]string) (interface{}, error) {
	db = form.whereSafe(db, parallel)
	if form.safeCounter == 0 {
		return db, errors.New("forbid no condition update")
	}
	return db, nil
}

func (form *FormBody) QueryCustom(f func(form *FormBody) interface{}) interface{} {
	return f(form)
}

// Create 创建model
func (form *FormBody) Create(model interface{}, db interface{}, parallel map[string]string) (interface{}, error) {
	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("CREATE", form.DriverName(parallel), "")
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
func (form *FormBody) CreateWithValid(model interface{}, db interface{}, parallel map[string]string, valid func(model interface{}) error) (interface{}, error) {
	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	err = valid(model)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("CREATE", form.DriverName(parallel), "")
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
func (form *FormBody) Update(model interface{}, db interface{}, parallel map[string]string) (interface{}, error) {

	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	// TODO 更新model必须提供有效的更新条件
	db, err = form.QuerySafe(db, parallel)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("UPDATES", form.DriverName(parallel), "")
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
func (form *FormBody) UpdateWithValid(model interface{}, db interface{}, parallel map[string]string, valid func(model interface{}) error) (interface{}, error) {

	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	err = valid(model)
	if err != nil {
		return nil, err
	}

	db, err = form.QuerySafe(db, parallel)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("UPDATES", form.DriverName(parallel), "")
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
func (form *FormBody) UpdateMapWithValid(model interface{}, db interface{}, parallel map[string]string, valid func(model interface{}) (error, map[string]interface{})) (interface{}, error) {

	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	err, m := valid(model)
	if err != nil {
		return nil, err
	}

	db, err = form.QuerySafe(db, parallel)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("UPDATES", form.DriverName(parallel), "")
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
func (form *FormBody) Save(model interface{}, db interface{}, parallel map[string]string) (interface{}, error) {

	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("SAVE", form.DriverName(parallel), "")
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
func (form *FormBody) SaveWithValid(model interface{}, db interface{}, parallel map[string]string, valid func(model interface{}) error) (interface{}, error) {

	err := form.Unmarshal(model)
	if err != nil {
		return nil, err
	}

	err = valid(model)
	if err != nil {
		return nil, err
	}

	exec := GetExecute("SAVE", form.DriverName(parallel), "")
	if exec == nil {
		return nil, errors.New("execution function for 'SAVE' not found")
	}

	if e := exec(db, "", model); e == nil {
		return model, nil
	} else {
		return nil, e.(error)
	}

}
