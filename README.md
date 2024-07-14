# gocrud

对常用的`create`、`read`、`update`、`delete`、`sort`进行封装，可实现具体`execute`接口进行扩展。

## 相关实现库

- [gocrud-gorm](https://github.com/kordar/gocrud-gorm)

## 常用操作符

| 名称   | 描述      | 必须 |
|------|---------|----|
| =    | 等于      | Y  |
| EQ   | 等于      | Y  |
| !=   | 不等于     | Y  |
| <>   | 不等于     | Y  |
| NEQ  | 不等于     | Y  |
| <    | 小于      | Y  |
| LT   | 小于      | Y  |
| <=   | 小于等于    | Y  |
| LE   | 小于等于    | Y  |
| &gt; | 大于      | Y  |
| GT   | 大于      | Y  |
| &gt;=   | 大于等于    | Y  |
| GE   | 大于等于    | Y  |
| IN   | 属于      | Y  |
| NOTIN   | 不属于     | Y  |
| LIKE   | 模糊查询    | Y  |
| NOTLIKE   | 模糊查询非   | Y  |
| LIKELEFT   | 模糊查询（左） | Y  |
| LIKERIGHT   | 模糊查询（右） | Y  |
| BETWEEN   | 范围      | Y  |
| NOTBETWEEN   | 范围外     | Y  |
| ISNULL   | 是否空     | Y  |
| ISNOTNULL   | 是否不为空   | Y  |
| ASC   | 升序      | Y  |
| DESC   | 降序      | Y  |
| SAVE   | 保存      | Y  |
| SETVAL   | 设置值     | Y  |
| UPDATE   | 更新值     | Y  |
| UPDATES   | 更新对象    | Y  |
| CREATE   | 创建      | Y  |
| PAGE   | 分页      | Y  |
| DELETE   | 删除      | Y  |

