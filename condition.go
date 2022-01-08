package gocrud

import (
	"fmt"
	"strings"
)

import (
	"github.com/kordar/goutil"
	"gorm.io/gorm"
)

var conditions = map[string]operator{
	"=":          &EQ{},
	"EQ":         &EQ{},
	"!=":         &NEQ{},
	"<>":         &NEQ{},
	"LT":         &LT{},
	"<":          &LT{},
	"LE":         &LE{},
	"<=":         &LE{},
	"GT":         &GT{},
	">":          &GT{},
	"GE":         &GE{},
	">=":         &GE{},
	"NEQ":        &NEQ{},
	"IN":         &IN{},
	"NOTIN":      &NOTIN{},
	"LIKE":       &LIKE{},
	"NOTLIKE":    &NOTLIKE{},
	"LIKELEFT":   &LIKELEFT{},
	"LIKERIGHT":  &LIKERIGHT{},
	"BETWEEN":    &BETWEEN{},
	"NOTBETWEEN": &NOTBETWEEN{},
	"ISNULL":     &ISNULL{},
	"ISNOTNULL":  &ISNOTNULL{},
}

type operator interface {
	execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB
}

type condition struct {
	Property    string      `json:"property" form:"property"`
	Key         string      `json:"key" form:"key"`
	Field       string      `json:"field" form:"field"`
	Value       interface{} `json:"value,omitempty" form:"value,omitempty"`
	Value2      interface{} `json:"value2,omitempty" form:"value2,omitempty"`
	Type        string      `json:"type" form:"type"`
	FilterEmpty bool        `json:"filter_empty" form:"filter_empty"`
}

func (c condition) Where(db *gorm.DB, parallel map[string]string) *gorm.DB {
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

	t := strings.ToUpper(c.Type)
	if t == "" {
		t = "EQ"
	}

	o := conditions[t]
	if o == nil {
		o = &EQ{}
	}

	return o.execute(db, field, c.Value, c.Value2)
}

// =
type EQ struct {
}

func (E *EQ) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s = ?", field), value)
}

// !=
type NEQ struct {
}

func (E *NEQ) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s <> ?", field), value)
}

// <
type LT struct {
}

func (E *LT) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s < ?", field), value)
}

// <=
type LE struct {
}

func (E *LE) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s <= ?", field), value)
}

// >
type GT struct {
}

func (E *GT) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s > ?", field), value)
}

// >=
type GE struct {
}

func (E *GE) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s >= ?", field), value)
}

// in
type IN struct {
}

func (E *IN) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s IN ?", field), value)
}

// not in
type NOTIN struct {
}

func (E *NOTIN) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s NOT IN ?", field), value)
}

// like
type LIKE struct {
}

func (E *LIKE) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	s := value.(string)
	return db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+s+"%")
}

// not like
type NOTLIKE struct {
}

func (E *NOTLIKE) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	s := value.(string)
	return db.Where(fmt.Sprintf("%s NOT LIKE %%?%%", field), "%"+s+"%")
}

// like left
type LIKELEFT struct {
}

func (E *LIKELEFT) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	s := value.(string)
	return db.Where(fmt.Sprintf("%s NOT LIKE ?%%", field), s+"%")
}

// like right
type LIKERIGHT struct {
}

func (E *LIKERIGHT) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	s := value.(string)
	return db.Where(fmt.Sprintf("%s NOT LIKE %%?", field), "%"+s)
}

// between
type BETWEEN struct {
}

func (E *BETWEEN) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), value, value2[0])
}

// not between
type NOTBETWEEN struct {
}

func (E *NOTBETWEEN) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s NOT BETWEEN ? AND ?", field), value, value2[0])
}

// is null
type ISNULL struct {
}

func (E *ISNULL) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s IS NULL", field))
}

// is not null
type ISNOTNULL struct {
}

func (E *ISNOTNULL) execute(db *gorm.DB, field string, value interface{}, value2 ...interface{}) *gorm.DB {
	return db.Where(fmt.Sprintf("%s IS NOT NULL", field))
}
