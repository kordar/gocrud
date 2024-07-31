package gocrud

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

// Execute create,update,delete返回error
type Execute func(db interface{}, field string, value interface{}, value2 ...interface{}) interface{}

func AddExecute(name string, execute Execute, driver string) {
	key := GetNameWithDriver(name, driver)
	execs[key] = execute
}

func GetExecute(name string, driver string, defaultName string) Execute {
	key := GetNameWithDriver(name, driver)
	exec := execs[key]
	if exec == nil {
		name = GetNameWithDriver(defaultName, driver)
		return execs[name]
	}
	return exec
}
