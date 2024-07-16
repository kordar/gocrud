package gocrud

import "errors"

type EditorBody[T interface{}, C interface{}] struct {
	Conditions  []Condition `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Editors     []Editor    `json:"editor,omitempty" form:"editor,omitempty"`
	safeCounter int         // 防止无条件更新
	*CommonBody[C]
}

func NewEditorBody[T interface{}, C interface{}](driver string, ctx C) EditorBody[T, C] {
	return EditorBody[T, C]{
		CommonBody: NewCommonBody[C](driver, ctx),
		Conditions: make([]Condition, 0),
	}
}

func (form *EditorBody[T, C]) where(db T, parallel map[string]string) T {
	parallel = form.LoadDriverName(parallel)
	for _, c := range form.Conditions {
		db = c.Where(db, parallel).(T)
	}
	return db
}

func (form *EditorBody[T, C]) whereSafe(db T, parallel map[string]string) T {
	parallel = form.LoadDriverName(parallel)
	for _, c := range form.Conditions {
		db2, flag := c.WhereSafe(db, parallel)
		db = db2.(T)
		if flag == true {
			form.safeCounter++
		}
	}
	return db
}

// Query 条件查询
func (form *EditorBody[T, C]) Query(db T, parallel map[string]string) interface{} {
	return form.where(db, parallel)
}

// QuerySafe 防止空条件更新
func (form *EditorBody[T, C]) QuerySafe(db T, parallel map[string]string) (interface{}, error) {
	db = form.whereSafe(db, parallel)
	if form.safeCounter == 0 {
		return db, errors.New("forbid no condition edit")
	}
	return db, nil
}

// UpdateData the data for update
func (form *EditorBody[T, C]) UpdateData(parallel map[string]string) map[string]interface{} {
	data := map[string]interface{}{}
	for _, exec := range form.Editors {
		k, v := exec.Param(parallel)
		data[k] = v
	}
	return data
}

// Updates update model object
func (form *EditorBody[T, C]) Updates(model interface{}, db T, parallel map[string]string) error {
	newDb, err := form.QuerySafe(db, parallel)
	if err != nil {
		return err
	}

	exec := GetExecute("UPDATES", form.DriverName(parallel), "")
	if exec == nil {
		return errors.New("execution function for 'UPDATES' not found")
	}

	data := form.UpdateData(parallel)
	if e := exec(newDb, "", model, data); e == nil {
		return nil
	} else {
		return e.(error)
	}

}

// QueryCustom 自定义条件查询
func (form *EditorBody[T, C]) QueryCustom(f func(form *EditorBody[T, C]) T) T {
	return f(form)
}
