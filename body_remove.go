package gocrud

import (
	"context"
	"errors"
	"github.com/kordar/gormext"
)

type RemoveBody struct {
	Data        map[string]interface{} `json:"data,omitempty" form:"data,omitempty"`             // 数据
	Conditions  []Condition            `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	safeCounter gormext.StrInt         // 防止无条件更新
	*CommonBody
}

func NewRemoveBody(driver string, ctx context.Context) RemoveBody {
	return RemoveBody{
		Data:       make(map[string]interface{}),
		Conditions: make([]Condition, 0),
		CommonBody: NewCommonBody(driver, ctx),
	}
}

func (remove *RemoveBody) where(db interface{}, params map[string]string) interface{} {
	params = remove.LoadDriverName(params)
	for _, exec := range remove.Conditions {
		db = exec.Where(db, params)
		remove.safeCounter++
	}
	return db
}

func (remove *RemoveBody) QuerySafe(db interface{}, params map[string]string) (interface{}, error) {
	db = remove.where(db, params)
	if remove.safeCounter == 0 {
		return db, errors.New("forbid no condition remove")
	}
	return db, nil
}

func (remove *RemoveBody) Query(db interface{}, params map[string]string) interface{} {
	return remove.where(db, params)
}

func (remove *RemoveBody) Delete(model interface{}, db interface{}, params map[string]string) error {
	if tx, err := remove.QuerySafe(db, nil); err != nil {
		return err
	} else {
		exec := GetExecute("DELETE", remove.DriverName(params), "")
		if exec == nil {
			return errors.New("execution function for 'DELETE' not found")
		}

		if e := exec(tx, "", model); e == nil {
			return nil
		} else {
			return e.(error)
		}
	}
}

func (remove *RemoveBody) QueryCustom(f func(form *RemoveBody) interface{}) interface{} {
	return f(remove)
}
