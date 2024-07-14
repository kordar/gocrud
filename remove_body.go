package gocrud

import (
	"errors"
)

type RemoveBody struct {
	Data        map[string]interface{} `json:"data,omitempty" form:"data,omitempty"`             // 数据
	Conditions  []condition            `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	safeCounter int                    // 防止无条件更新
	commonBody
}

func NewRemoveBody(driver string) RemoveBody {
	return RemoveBody{
		Data:       make(map[string]interface{}),
		Conditions: make([]condition, 0),
		commonBody: commonBody{Driver: driver},
	}
}

func (remove *RemoveBody) where(db interface{}, parallel map[string]string) interface{} {
	for _, exec := range remove.Conditions {
		db = exec.Where(db, parallel)
		remove.safeCounter++
	}
	return db
}

func (remove *RemoveBody) QuerySafe(db interface{}, parallel map[string]string) (interface{}, error) {
	db = remove.where(db, parallel)
	if remove.safeCounter == 0 {
		return db, errors.New("forbid no condition remove")
	}
	return db, nil
}

func (remove *RemoveBody) Query(db interface{}, parallel map[string]string) interface{} {
	return remove.where(db, parallel)
}

func (remove *RemoveBody) QueryCustom(f func(form *RemoveBody) interface{}) interface{} {
	return f(remove)
}

func (remove *RemoveBody) Delete(model interface{}, db interface{}, parallel map[string]string) error {
	if tx, err := remove.QuerySafe(db, nil); err != nil {
		return err
	} else {
		exec := GetExecute("DELETE", remove.GetDriver(parallel), "")
		if e := exec(tx, "", model); e == nil {
			return nil
		} else {
			return e.(error)
		}
	}
}
