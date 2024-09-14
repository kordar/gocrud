package gocrud

import "github.com/kordar/gologger"

type Condition struct {
	Property    string      `json:"property" form:"property"`
	Key         string      `json:"key" form:"key"`
	Field       string      `json:"field" form:"field"`
	Value       interface{} `json:"value,omitempty" form:"value,omitempty"`
	Value2      interface{} `json:"value2,omitempty" form:"value2,omitempty"`
	Type        string      `json:"type" form:"type"`
	FilterEmpty bool        `json:"filter_empty" form:"filter_empty"`
}

func (c Condition) WhereSafe(db interface{}, params map[string]string) (interface{}, bool) {

	if c.Value == nil {
		return db, false
	}

	// 过滤空值
	if c.FilterEmpty && (c.Value == "" || c.Value == 0) {
		return db, false
	}

	/**
	 * 获取属性
	 * {"property": "属性", "key": "键值", "field": "字段值"}
	 * 存在属性值以属性值为准，否则将key值计算驼峰设置为属性值
	 */
	realField := GetRealFiled(c.Field, c.Property, c.Key, params)

	if realField == "" {
		return db, false
	}

	exec := GetExecute(c.Type, params["driver"], "EQ")
	if exec == nil {
		logger.Warnf("[gocrud] execution function for '%s' not found", c.Type)
		return db, false
	}

	return exec(db, realField, c.Value, c.Value2), true
}

func (c Condition) Where(db interface{}, params map[string]string) interface{} {

	if c.Value == nil {
		return db
	}

	// 过滤空值
	if c.FilterEmpty && (c.Value == "" || c.Value == 0) {
		return db
	}

	/**
	 * 获取属性
	 * {"property": "属性", "key": "键值", "field": "字段值"}
	 * 存在属性值以属性值为准，否则将key值计算驼峰设置为属性值
	 */
	realField := GetRealFiled(c.Field, c.Property, c.Key, params)

	if realField == "" {
		return db
	}

	exec := GetExecute(c.Type, params["driver"], "EQ")
	if exec == nil {
		logger.Warnf("[gocrud] execution function for '%s' not found", c.Type)
		return db
	}

	return exec(db, realField, c.Value, c.Value2)
}
