package gocrud

import (
	"encoding/json"
	"errors"
)

type FormBody[T interface{}, C interface{}] struct {
	Conditions  []Condition `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Object      interface{} `json:"object,omitempty" form:"object,omitempty"`
	safeCounter int         // 防止无条件更新
	*CommonBody[C]
}

func NewFormBody[T interface{}, C interface{}](driver string, ctx C) FormBody[T, C] {
	return FormBody[T, C]{
		CommonBody: NewCommonBody[C](driver, ctx),
		Conditions: make([]Condition, 0),
	}
}

func (form *FormBody[T, C]) Unmarshal(target interface{}) error {
	if marshal, err := json.Marshal(form.Object); err != nil {
		return err
	} else {
		return json.Unmarshal(marshal, target)
	}
}

func (form *FormBody[T, C]) where(db T, parallel map[string]string) T {
	parallel = form.LoadDriverName(parallel)
	if form.Conditions != nil {
		for _, exec := range form.Conditions {
			db = exec.Where(db, parallel).(T)
		}
	}
	return db
}

func (form *FormBody[T, C]) whereSafe(db T, parallel map[string]string) T {
	parallel = form.LoadDriverName(parallel)
	if form.Conditions != nil {
		for _, exec := range form.Conditions {
			db2, flag := exec.WhereSafe(db, parallel)
			db = db2.(T)
			if flag == true {
				form.safeCounter++
			}
		}
	}
	return db
}

func (form *FormBody[T, C]) Query(db T, parallel map[string]string) T {
	return form.where(db, parallel)
}

func (form *FormBody[T, C]) QuerySafe(db T, parallel map[string]string) (T, error) {
	db = form.whereSafe(db, parallel)
	if form.safeCounter == 0 {
		return db, errors.New("forbid no condition update")
	}
	return db, nil
}

func (form *FormBody[T, C]) QueryCustom(f func(form *FormBody[T, C]) T) T {
	return f(form)
}

// Create 创建model
func (form *FormBody[T, C]) Create(model interface{}, db T, parallel map[string]string) (interface{}, error) {
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
func (form *FormBody[T, C]) CreateWithValid(model interface{}, db T, parallel map[string]string, valid func(model interface{}) error) (interface{}, error) {
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
func (form *FormBody[T, C]) Update(model interface{}, db T, parallel map[string]string) (interface{}, error) {

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
func (form *FormBody[T, C]) UpdateWithValid(model interface{}, db T, parallel map[string]string, valid func(model interface{}) error) (interface{}, error) {

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
func (form *FormBody[T, C]) UpdateMapWithValid(model interface{}, db T, parallel map[string]string, valid func(model interface{}) (error, map[string]interface{})) (interface{}, error) {

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

// SaveWithValid 更新并且校验提交参数
func (form *FormBody[T, C]) Save(model interface{}, db T, parallel map[string]string) (interface{}, error) {

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
func (form *FormBody[T, C]) SaveWithValid(model interface{}, db T, parallel map[string]string, valid func(model interface{}) error) (interface{}, error) {

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
