package gocrud

import "errors"

type RemoveBody[T interface{}, C interface{}] struct {
	Data        map[string]interface{} `json:"data,omitempty" form:"data,omitempty"`             // 数据
	Conditions  []Condition            `json:"conditions,omitempty" form:"conditions,omitempty"` // 条件
	safeCounter int                    // 防止无条件更新
	*CommonBody[C]
}

func NewRemoveBody[T interface{}, C interface{}](driver string, ctx C) RemoveBody[T, C] {
	return RemoveBody[T, C]{
		Data:       make(map[string]interface{}),
		Conditions: make([]Condition, 0),
		CommonBody: NewCommonBody[C](driver, ctx),
	}
}

func (remove *RemoveBody[T, C]) where(db T, parallel map[string]string) T {
	parallel = remove.LoadDriverName(parallel)
	for _, exec := range remove.Conditions {
		db = exec.Where(db, parallel).(T)
		remove.safeCounter++
	}
	return db
}

func (remove *RemoveBody[T, C]) QuerySafe(db T, parallel map[string]string) (T, error) {
	db = remove.where(db, parallel)
	if remove.safeCounter == 0 {
		return db, errors.New("forbid no condition remove")
	}
	return db, nil
}

func (remove *RemoveBody[T, C]) Query(db T, parallel map[string]string) T {
	return remove.where(db, parallel)
}

func (remove *RemoveBody[T, C]) Delete(model interface{}, db T, parallel map[string]string) error {
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

func (remove *RemoveBody[T, C]) QueryCustom(f func(form *RemoveBody[T, C]) T) T {
	return f(remove)
}
