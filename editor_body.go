package gocrud

import (
	"errors"
)

type EditorBody struct {
	Conditions  []condition `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Editors     []editor    `json:"editor,omitempty" form:"editor,omitempty"`
	safeCounter int         // 防止无条件更新
	commonBody
}

func NewEditorBody(driver string) EditorBody {
	return EditorBody{
		commonBody: commonBody{Driver: driver},
		Conditions: make([]condition, 0),
	}
}

func (form *EditorBody) where(db interface{}, parallel map[string]string) interface{} {
	for _, c := range form.Conditions {
		db = c.Where(db, parallel)
	}
	return db
}

func (form *EditorBody) whereSafe(db interface{}, parallel map[string]string) interface{} {
	for _, c := range form.Conditions {
		flag := false
		db, flag = c.WhereSafe(db, parallel)
		if flag == true {
			form.safeCounter++
		}
	}
	return db
}

// UpdateData the data for update
func (form *EditorBody) UpdateData(parallel map[string]string) map[string]interface{} {
	data := map[string]interface{}{}
	for _, exec := range form.Editors {
		k, v := exec.Param(parallel)
		data[k] = v
	}
	return data
}

// Updates update model object
func (form *EditorBody) Updates(db interface{}, parallel map[string]string) error {
	newDb, err := form.QuerySafe(db, parallel)
	if err != nil {
		return err
	}

	exec := GetExecute("UPDATES", form.GetDriver(parallel), "")
	data := form.UpdateData(parallel)
	if e := exec(newDb, "", data); e == nil {
		return nil
	} else {
		return e.(error)
	}

}

// QuerySafe 防止空条件更新
func (form *EditorBody) QuerySafe(db interface{}, parallel map[string]string) (interface{}, error) {
	db = form.whereSafe(db, parallel)
	if form.safeCounter == 0 {
		return db, errors.New("forbid no condition edit")
	}
	return db, nil
}

// Query 条件查询
func (form *EditorBody) Query(db interface{}, parallel map[string]string) interface{} {
	return form.where(db, parallel)
}

// QueryCustom 自定义条件查询
func (form *EditorBody) QueryCustom(
	db interface{},
	parallel map[string]string,
	fun func(form *EditorBody, db interface{}, parallel map[string]string) interface{},
) interface{} {
	return fun(form, db, parallel)
}
