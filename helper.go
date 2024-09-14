package gocrud

import (
	"context"
	"fmt"
	"github.com/kordar/goutil"
	"strings"
)

var MessageFn = func(ctx context.Context, message string) string {
	return message
}

func GetNameWithDriver(name string, driver string) string {
	if driver == "" {
		return strings.ToUpper(name)
	} else {
		return fmt.Sprintf("%s:%s", strings.ToUpper(driver), strings.ToUpper(name))
	}
}

// GetRealFieldValue 通过__Field 设置隐藏参数的真实名称
func GetRealFieldValue(key string, parallel map[string]string) string {
	field := GetField(key)
	return parallel[field]
}

func GetField(key string) string {
	return fmt.Sprintf("__%s", key)
}

func GetRealFiled(field string, property string, key string, parallel map[string]string) string {
	// field 不为空，则判断是否存在field映射
	realField := ""

	key = strings.Trim(key, " ")
	if key != "" {
		realField = goutil.SnakeString(key)
	}

	property = strings.Trim(property, " ")
	if property != "" {
		realField = goutil.SnakeString(property)
	}

	field = strings.Trim(field, " ")
	if field != "" {
		realField = field
	}

	if realField != "" {
		targetField := GetRealFieldValue(realField, parallel)
		if targetField != "" {
			return targetField
		}
	}

	return realField
}
