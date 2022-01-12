package gocrud

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type commonBody struct {
	Ctx *gin.Context
}

func (common *commonBody) Lang() string {
	return common.Ctx.DefaultQuery("locale", "en")
}

type FormBody struct {
	Conditions  []*condition `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	Object      interface{}  `json:"object,omitempty" form:"object,omitempty"`
	safeCounter int          // 防止无条件更新
	commonBody
}

func NewFormBody(ctx *gin.Context) FormBody {
	return FormBody{
		commonBody: commonBody{Ctx: ctx},
		Conditions: make([]*condition, 0),
	}
}

func (form *FormBody) GetObject(target interface{}) error {
	if marshal, err := json.Marshal(form.Object); err != nil {
		return err
	} else {
		return json.Unmarshal(marshal, target)
	}
}

func (form *FormBody) where(db *gorm.DB, parallel map[string]string) *gorm.DB {
	for _, condition := range form.Conditions {
		db = condition.Where(db, parallel)
	}
	return db
}

func (form *FormBody) whereSafe(db *gorm.DB, parallel map[string]string) *gorm.DB {
	for _, condition := range form.Conditions {
		flag := false
		db, flag = condition.WhereSafe(db, parallel)
		if flag == true {
			form.safeCounter++
		}
	}
	return db
}

func (form *FormBody) QuerySafe(db *gorm.DB, parallel map[string]string) (*gorm.DB, error) {
	db = form.whereSafe(db, parallel)
	if form.safeCounter == 0 {
		return db, errors.New("forbid no condition update")
	}
	return db, nil
}

func (form *FormBody) Query(db *gorm.DB, parallel map[string]string) *gorm.DB {
	return form.where(db, parallel)
}

func (form *FormBody) QueryCustom(db *gorm.DB, parallel map[string]string, fun func(form *FormBody, db *gorm.DB, parallel map[string]string) *gorm.DB) *gorm.DB {
	return fun(form, db, parallel)
}
