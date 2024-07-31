package gocrud

import "errors"

type RemoveBody struct {
	Data        map[string]interface{} `json:"data,omitempty" form:"data,omitempty"`             // 数据
	Conditions  []Condition            `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	safeCounter int                    // 防止无条件更新
	*CommonBody
}

func NewRemoveBody(driver string, ctx interface{}) RemoveBody {
	return RemoveBody{
		Data:       make(map[string]interface{}),
		Conditions: make([]Condition, 0),
		CommonBody: NewCommonBody(driver, ctx),
	}
}

func (remove *RemoveBody) where(db interface{}, parallel map[string]string) interface{} {
	parallel = remove.LoadDriverName(parallel)
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

func (remove *RemoveBody) Delete(model interface{}, db interface{}, parallel map[string]string) error {
	if tx, err := remove.QuerySafe(db, nil); err != nil {
		return err
	} else {
		exec := GetExecute("DELETE", remove.DriverName(parallel), "")
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
