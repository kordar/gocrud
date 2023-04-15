package gocrud

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EditorBody struct {
	Conditions  []*condition `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Editors     []*editor    `json:"editor,omitempty" form:"editor,omitempty"`
	safeCounter int          // 防止无条件更新
	commonBody
}

func NewEditorBody(ctx *gin.Context) EditorBody {
	return EditorBody{
		commonBody: commonBody{Ctx: ctx},
		Conditions: make([]*condition, 0),
	}
}

func (form *EditorBody) where(db *gorm.DB, parallel map[string]string) *gorm.DB {
	for _, condition := range form.Conditions {
		db = condition.Where(db, parallel)
	}
	return db
}

func (form *EditorBody) whereSafe(db *gorm.DB, parallel map[string]string) *gorm.DB {
	for _, condition := range form.Conditions {
		flag := false
		db, flag = condition.WhereSafe(db, parallel)
		if flag == true {
			form.safeCounter++
		}
	}
	return db
}

// UpdateData the data for update
func (form *EditorBody) UpdateData(parallel map[string]string) map[string]interface{} {
	data := map[string]interface{}{}
	for _, editor := range form.Editors {
		k, v := editor.Param(parallel)
		data[k] = v
	}
	return data
}

// Updates update model object
func (form *EditorBody) Updates(model interface{}, db *gorm.DB, parallel map[string]string) error {
	if db, err := form.QuerySafe(db, parallel); err != nil {
		return err
	} else {
		data := form.UpdateData(parallel)
		db := db.Model(model).UpdateColumns(data)
		return db.Error
	}
}

// QuerySafe 防止空条件更新
func (form *EditorBody) QuerySafe(db *gorm.DB, parallel map[string]string) (*gorm.DB, error) {
	db = form.whereSafe(db, parallel)
	if form.safeCounter == 0 {
		return db, errors.New("forbid no condition edit")
	}
	return db, nil
}

// Query 条件查询
func (form *EditorBody) Query(db *gorm.DB, parallel map[string]string) *gorm.DB {
	return form.where(db, parallel)
}

// QueryCustom 自定义条件查询
func (form *EditorBody) QueryCustom(db *gorm.DB, parallel map[string]string, fun func(form *EditorBody, db *gorm.DB, parallel map[string]string) *gorm.DB) *gorm.DB {
	return fun(form, db, parallel)
}
