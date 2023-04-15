package gocrud

import (
	"fmt"
	"github.com/kordar/goutil/helper"
	"gorm.io/gorm"
	"strings"
)

var sorts = map[string]sortOperator{
	"ASC":  &ASC{},
	"DESC": &DESC{},
}

type sortOperator interface {
	execute(db *gorm.DB, field string) *gorm.DB
}

type sort struct {
	Property string `json:"property" form:"property"`
	Key      string `json:"key" form:"key"`
	Field    string `json:"field" form:"field"`
	Type     string `json:"type" form:"type"`
}

func (c sort) Order(db *gorm.DB, parallel map[string]string) *gorm.DB {
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
		field = helper.SnakeString(property)
	}

	t := strings.ToUpper(c.Type)
	if t == "" {
		t = "ASC"
	}

	o := sorts[t]
	if o == nil {
		o = &ASC{}
	}

	return o.execute(db, field)
}

type ASC struct {
}

func (A *ASC) execute(db *gorm.DB, field string) *gorm.DB {
	return db.Order(fmt.Sprintf("%s ASC", field))
}

type DESC struct {
}

func (A *DESC) execute(db *gorm.DB, field string) *gorm.DB {
	return db.Order(fmt.Sprintf("%s DESC", field))
}
