# gocrud

对常用的`create`、`read`、`update`、`delete`、`sort`进行封装，可实现具体`execute`接口进行扩展。

## 相关实现库

- [gocrud-gorm](https://github.com/kordar/gocrud-gorm)

## 常用操作符

| 名称   | 描述      | 返回值   |
|------|---------|-------|
| =    | 等于      | 句柄    |
| EQ   | 等于      | 句柄    |
| !=   | 不等于     | 句柄    |
| <>   | 不等于     | 句柄    |
| NEQ  | 不等于     | 句柄    |
| <    | 小于      | 句柄    |
| LT   | 小于      | 句柄    |
| <=   | 小于等于    | 句柄    |
| LE   | 小于等于    | 句柄    |
| &gt; | 大于      | 句柄    |
| GT   | 大于      | 句柄    |
| &gt;=   | 大于等于    | 句柄    |
| GE   | 大于等于    | 句柄    |
| IN   | 属于      | 句柄    |
| NOTIN   | 不属于     | 句柄    |
| LIKE   | 模糊查询    | 句柄    |
| NOTLIKE   | 模糊查询非   | 句柄    |
| LIKELEFT   | 模糊查询（左） | 句柄    |
| LIKERIGHT   | 模糊查询（右） | 句柄    |
| BETWEEN   | 范围      | 句柄    |
| NOTBETWEEN   | 范围外     | 句柄    |
| ISNULL   | 是否空     | 句柄    |
| ISNOTNULL   | 是否不为空   | 句柄    |
| ASC   | 升序      | 句柄    |
| DESC   | 降序      | 句柄    |
| SAVE   | 保存      | error |
| SETVAL   | 设置值     | error |
| UPDATE   | 更新值     | error |
| UPDATES   | 更新对象    | error |
| CREATE   | 创建      | error |
| PAGE   | 分页      | 句柄    |
| DELETE   | 删除      | error |

## 使用方式

```go
// 设置多语言
func SetLangFunc(f func() string)
// 实现resource接口
type ResourceService[T interface{}, C interface{}] interface {
    ResourceName() string
    Search(body SearchBody[T, C]) SearchVO
    SearchOne(body SearchBody[T, C]) SearchOneVO
    Delete(body RemoveBody[T, C]) error
    Add(body FormBody[T, C]) (interface{}, error)
    Update(body FormBody[T, C]) (interface{}, error)
    Edit(body EditorBody[T, C]) error
    Configs() map[string]interface{}
    DriverName() string
}
// 初始化resource service管理容器
var Manager = gocrud.NewResourceManager()
// 注入实现类
Manager.AddResourceService(ResourceService)
```
