# gocrud

> 一个 **基于条件描述 + Driver 执行器的通用 CRUD 抽象层**
> 用于快速构建 **可配置 / 可扩展 / 多 Driver** 的资源型 CRUD 服务

---

## ✨ 设计目标

* 🚀 用 **JSON / Form 描述 CRUD 行为**
* 🔌 解耦业务逻辑与 ORM / Driver（如 GORM）
* 🧩 条件 / 排序 / 编辑 / 执行器高度可扩展
* 🛡️ 内置 **无条件更新 / 删除保护**
* 🧠 统一 CRUD 编排（Search / Create / Update / Delete / Edit）
* 📦 适配 REST / 管理后台 / 配置化接口

---

## 📦 安装

```bash
go get github.com/kordar/gocrud
```

---

## 🧱 核心概念总览

```tex
Request Body
   │
   ▼
SearchBody / FormBody / EditorBody / RemoveBody
   │
   ▼
Condition / Sort / Editor
   │
   ▼
Execute（Driver 绑定）
   │
   ▼
ORM / DB / Storage
```

---

## 🔧 CommonBody（公共上下文）

```go
type CommonBody struct {
	ctx    context.Context
	driver string
}
```

### 功能

* 统一管理：

  * `context.Context`
  * 当前 Driver
* 支持 **参数中动态切换 driver**

### 使用

```go
body := NewCommonBody("gorm", ctx)
driver := body.DriverName(params)
```

---

## 🔍 SearchBody（查询）

用于 **列表 / 单条查询**

```go
type SearchBody struct {
	Page       int
	PageSize   int
	Conditions []Condition
	Sorts      []Sort
}
```

### 示例

```go
search := NewSearchBody("gorm", ctx)
search.Conditions = []Condition{
	{Key: "name", Type: "LIKE", Value: "Tom"},
}
search.Sorts = []Sort{
	{Key: "created_at", Type: "DESC"},
}
```

### 执行流程

```go
db = search.Query(db, params)
db = search.Paginate(db, params)
```

---

## 📝 Condition（查询条件）

```go
type Condition struct {
	Property    string
	Key         string
	Field       string
	Value       interface{}
	Value2      interface{}
	Type        string
	FilterEmpty bool
}
```

### 字段解析优先级

```
Field > Property > Key
```

最终统一映射为 **数据库字段名**

### 示例

```json
{
  "key": "name",
  "type": "LIKE",
  "value": "Tom"
}
```

---

## 🔃 Sort（排序）

```go
type Sort struct {
	Key   string
	Field string
	Type  string
}
```

### 示例

```json
{
  "key": "created_at",
  "type": "DESC"
}
```

---

## 🧾 FormBody（Create / Update / Save）

用于 **模型级 CRUD**

```go
type FormBody struct {
	Object     interface{}
	Conditions []Condition
}
```

### Create

```go
form := NewFormBody("gorm", ctx)
form.Object = reqBody

result, err := form.Create(&User{}, db, nil)
```

---

### Update（安全）

* **必须有有效 Condition**
* 否则直接拒绝

```go
form.Conditions = []Condition{
	{Key: "id", Type: "EQ", Value: 1},
}

form.Update(&User{}, db, nil)
```

---

### Save（不校验条件）

```go
form.Save(&User{}, db, nil)
```

---

## ✏️ EditorBody（字段级更新）

用于 **部分字段更新（Patch / 批量更新）**

```go
type EditorBody struct {
	Conditions []Condition
	Editors    []Editor
}
```

### Editor 示例

```go
Editor{
	Key:   "status",
	Type: "SETVAL",
	Value: 1,
}
```

### 执行

```go
editor := NewEditorBody("gorm", ctx)
editor.Editors = []Editor{...}
editor.Conditions = []Condition{...}

editor.Updates(&User{}, db, nil)
```

---

## 🗑 RemoveBody（删除）

```go
type RemoveBody struct {
	Conditions []Condition
}
```

### 特性

* **无条件删除直接拒绝**
* 强制安全校验

```go
remove := NewRemoveBody("gorm", ctx)
remove.Conditions = []Condition{
	{Key: "id", Type: "EQ", Value: 1},
}
remove.Delete(&User{}, db, nil)
```

---

## ⚙️ Execute（执行器核心）

```go
type Execute func(
	db interface{},
	field string,
	value interface{},
	value2 ...interface{},
) interface{}
```

### 注册 Execute

```go
AddExecute("EQ", func(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	return db.Where(field+" = ?", value)
}, "gorm")
```

### Driver 绑定规则

```
GORM:EQ
MYSQL:LIKE
```

---

## 🔌 Driver 解耦机制

```go
GetNameWithDriver("EQ", "gorm")
// => GORM:EQ
```

支持：

* 同一 CRUD 逻辑
* 不同 ORM / 存储引擎

---

## 🧠 ResourceService（资源抽象）

```go
type ResourceService interface {
	Search(SearchBody) SearchVO
	Create(FormBody) (interface{}, error)
	Update(FormBody) (interface{}, error)
	Remove(RemoveBody) error
}
```

### 示例

```go
type UserService struct {
	CommonResourceService
}
```

---

## 📦 ResourceManager（统一入口）

```go
mgr := NewResourceManager()
mgr.AddResourceService(userService)

mgr.Read("user", searchBody)
mgr.Create("user", formBody)
mgr.Delete("user", removeBody)
```

---

## 📤 返回结构

### SearchVO

```go
type SearchVO struct {
	Data  interface{}
	Count int64
}
```

### SearchOneVO

```go
type SearchOneVO struct {
	Info interface{}
}
```

---

## 🛡️ 安全机制

| 场景         | 处理       |
| ---------- | -------- |
| 无条件 Update | ❌ 拒绝     |
| 无条件 Delete | ❌ 拒绝     |
| 空值过滤       | 可配置      |
| Driver 未注册 | fallback |

---

## 🧪 适用场景

* 管理后台 CRUD
* 通用 REST CRUD API
* 配置驱动的查询系统
* 多数据库 / 多 ORM 项目
* 低代码 / 表单引擎后端

---

## 📄 License

MIT

