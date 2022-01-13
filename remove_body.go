package gocrud

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RemoveBody struct {
	Ctx         *gin.Context
	Data        map[string]interface{} `json:"data,omitempty" form:"data,omitempty"`             // 数据
	Conditions  []*condition           `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	safeCounter int                    // 防止无条件更新
	commonBody
}

func NewRemoveBody(ctx *gin.Context) RemoveBody {
	return RemoveBody{
		Data:       make(map[string]interface{}),
		Conditions: make([]*condition, 0),
		commonBody: commonBody{Ctx: ctx},
	}
}

func (remove *RemoveBody) where(db *gorm.DB, parallel map[string]string) *gorm.DB {
	for _, condition := range remove.Conditions {
		db = condition.Where(db, parallel)
		remove.safeCounter++
	}
	return db
}

func (remove *RemoveBody) QuerySafe(db *gorm.DB, parallel map[string]string) (*gorm.DB, error) {
	db = remove.where(db, parallel)
	if remove.safeCounter == 0 {
		return db, errors.New("forbid no condition remove")
	}
	return db, nil
}

func (remove *RemoveBody) Query(db *gorm.DB, parallel map[string]string) *gorm.DB {
	return remove.where(db, parallel)
}

func (remove *RemoveBody) QueryCustom(db *gorm.DB, parallel map[string]string, fun func(form *RemoveBody, db *gorm.DB, parallel map[string]string) *gorm.DB) *gorm.DB {
	return fun(remove, db, parallel)
}

func (remove *RemoveBody) Delete(model interface{}, db *gorm.DB) error {
	if tx, err := remove.QuerySafe(db, nil); err != nil {
		return err
	} else {
		return tx.Delete(model).Error
	}
}
