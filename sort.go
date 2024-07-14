package gocrud

import (
	"github.com/kordar/goutil"
	"strings"
)

type sort struct {
	Property string `json:"property" form:"property"`
	Key      string `json:"key" form:"key"`
	Field    string `json:"field" form:"field"`
	Type     string `json:"type" form:"type"`
}

func (c sort) Order(db interface{}, parallel map[string]string) interface{} {
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

	exec := GetExecute(c.Type, parallel["driver"], "ASC")
	return exec(db, field, "")
}
