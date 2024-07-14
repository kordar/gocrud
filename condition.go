package gocrud

import (
	"github.com/kordar/goutil"
	"strings"
)

type condition struct {
	Property    string      `json:"property" form:"property"`
	Key         string      `json:"key" form:"key"`
	Field       string      `json:"field" form:"field"`
	Value       interface{} `json:"value,omitempty" form:"value,omitempty"`
	Value2      interface{} `json:"value2,omitempty" form:"value2,omitempty"`
	Type        string      `json:"type" form:"type"`
	FilterEmpty bool        `json:"filter_empty" form:"filter_empty"`
}

func (c condition) WhereSafe(db interface{}, parallel map[string]string) (interface{}, bool) {
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
	property := strings.Trim(c.Property, " ")
	if property == "" {
		property = strings.Trim(c.Key, " ")
	}

	if property == "" {
		return db, false
	}

	/**
	* 通过属性值获取字段类型
	 */
	var field string
	if parallel == nil {
		field = parallel[property]
	}

	// 属性值转为下划线
	if field == "" {
		field = goutil.SnakeString(property)
	}

	exec := GetExecute(c.Type, parallel["driver"], "EQ")
	return exec(db, field, c.Value, c.Value2), true
}

func (c condition) Where(db interface{}, parallel map[string]string) interface{} {
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
	property := strings.Trim(c.Property, " ")
	if property == "" {
		property = strings.Trim(c.Key, " ")
	}

	if property == "" {
		return db
	}

	/**
	* 通过属性值获取字段类型
	 */
	var field string
	if parallel == nil {
		field = parallel[property]
	}

	// 属性值转为下划线
	if field == "" {
		field = goutil.SnakeString(property)
	}

	exec := GetExecute(c.Type, parallel["driver"], "EQ")
	return exec(db, field, c.Value, c.Value2)
}
