package gocrud

import (
	"fmt"
	"strings"
)

type Operator = string

const (
	SAVE    Operator = "save"
	UPDATES Operator = "updates"
	CREATE  Operator = "create"
	DELETE  Operator = "delete"

	// PAGE 分页，value=offset, value2[0]=limit
	PAGE Operator = "page"
)

var execs = map[string]Execute{}

// Execute Create,update,delete返回error
type Execute func(db interface{}, field string, value interface{}, value2 ...interface{}) interface{}

func GetName(name string, driver string) string {
	if driver == "" {
		return strings.ToUpper(name)
	} else {
		return fmt.Sprintf("%s:%s", strings.ToUpper(driver), strings.ToUpper(name))
	}
}

func AddExecute(name string, execute Execute, driver string) {
	key := GetName(name, driver)
	execs[key] = execute
}

func GetExecute(name string, driver string, defaultName string) Execute {
	key := GetName(name, driver)
	exec := execs[key]
	if exec == nil {
		name = GetName(defaultName, driver)
		return execs[name]
	}
	return exec
}
