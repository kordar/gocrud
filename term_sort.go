package gocrud

import (
	logger "github.com/kordar/gologger"
)

type Sort struct {
	Property string `json:"property" form:"property"`
	Key      string `json:"key" form:"key"`
	Field    string `json:"field" form:"field"`
	Type     string `json:"type" form:"type"`
}

func (c Sort) Order(db interface{}, params map[string]string) interface{} {
	/**
	 * 获取属性
	 * {"property": "属性", "key": "键值", "field": "字段值"}
	 * 存在属性值以属性值为准，否则将key值计算驼峰设置为属性值
	 */
	realField := GetRealFiled(c.Field, c.Property, c.Key, params)

	exec := GetExecute(c.Type, params["driver"], "ASC")
	if exec == nil {
		logger.Warnf("[gocrud] execution function for '%s' not found", c.Type)
		return db
	}

	return exec(db, realField, "")
}
