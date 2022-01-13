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

func (form *FormBody) Query(db *gorm.DB, parallel map[string]string) *gorm.DB {
	return form.where(db, parallel)
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

func (form *FormBody) QueryCustom(db *gorm.DB, parallel map[string]string, fun func(form *FormBody, db *gorm.DB, parallel map[string]string) *gorm.DB) *gorm.DB {
	return fun(form, db, parallel)
}

// Create 创建model
func (form *FormBody) Create(model interface{}, db *gorm.DB, parallel map[string]string) (interface{}, error) {
	err := form.GetObject(model)
	if err != nil {
		return nil, err
	}

	err = db.Create(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

// CreateWithValid 创建并且校验提交参数
func (form *FormBody) CreateWithValid(model interface{}, db *gorm.DB, parallel map[string]string, valid func(ctx *gin.Context, model interface{}) error) (interface{}, error) {
	err := form.GetObject(model)
	if err != nil {
		return nil, err
	}

	err = valid(form.Ctx, model)
	if err != nil {
		return nil, err
	}

	err = db.Create(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

// Update 更新
func (form *FormBody) Update(model interface{}, db *gorm.DB, parallel map[string]string) (interface{}, error) {

	err := form.GetObject(model)
	if err != nil {
		return nil, err
	}

	db, err = form.QuerySafe(db, nil)
	if err != nil {
		return nil, err
	}

	err = db.Updates(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

// UpdateWithValid 更新并且校验提交参数
func (form *FormBody) UpdateWithValid(model interface{}, db *gorm.DB, parallel map[string]string, valid func(ctx *gin.Context, model interface{}) error) (interface{}, error) {

	err := form.GetObject(model)
	if err != nil {
		return nil, err
	}

	err = valid(form.Ctx, model)
	if err != nil {
		return nil, err
	}

	db, err = form.QuerySafe(db, nil)
	if err != nil {
		return nil, err
	}

	err = db.Updates(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}
