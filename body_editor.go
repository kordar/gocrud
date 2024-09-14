package gocrud

import (
	"context"
	"errors"
)

type EditorBody struct {
	Conditions  []Condition `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Editors     []Editor    `json:"editor,omitempty" form:"editor,omitempty"`
	safeCounter int         // 防止无条件更新
	*CommonBody
}

func NewEditorBody(driver string, ctx context.Context) EditorBody {
	return EditorBody{
		CommonBody: NewCommonBody(driver, ctx),
		Conditions: make([]Condition, 0),
	}
}

func (form *EditorBody) where(db interface{}, params map[string]string) interface{} {
	params = form.LoadDriverName(params)
	for _, c := range form.Conditions {
		db = c.Where(db, params)
	}
	return db
}

func (form *EditorBody) whereSafe(db interface{}, params map[string]string) interface{} {
	params = form.LoadDriverName(params)
	for _, c := range form.Conditions {
		db2, flag := c.WhereSafe(db, params)
		db = db2
		if flag == true {
			form.safeCounter++
		}
	}
	return db
}

// Query 条件查询
func (form *EditorBody) Query(db interface{}, params map[string]string) interface{} {
	return form.where(db, params)
}

// QuerySafe 防止空条件更新
func (form *EditorBody) QuerySafe(db interface{}, params map[string]string) (interface{}, error) {
	db = form.whereSafe(db, params)
	if form.safeCounter == 0 {
		return db, errors.New("forbid no condition edit")
	}
	return db, nil
}

// UpdateData the data for update
func (form *EditorBody) UpdateData(params map[string]string) map[string]interface{} {
	data := map[string]interface{}{}
	for _, exec := range form.Editors {
		k, v := exec.Param(params)
		data[k] = v
	}
	return data
}

// Updates update model object
func (form *EditorBody) Updates(model interface{}, db interface{}, params map[string]string) error {
	newDb, err := form.QuerySafe(db, params)
	if err != nil {
		return err
	}

	exec := GetExecute("UPDATES", form.DriverName(params), "")
	if exec == nil {
		return errors.New("execution function for 'UPDATES' not found")
	}

	data := form.UpdateData(params)
	if e := exec(newDb, "", model, data); e == nil {
		return nil
	} else {
		return e.(error)
	}

}

// QueryCustom 自定义条件查询
func (form *EditorBody) QueryCustom(f func(form *EditorBody) interface{}) interface{} {
	return f(form)
}
