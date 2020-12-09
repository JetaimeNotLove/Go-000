# 问题

`我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？`

# 答

对于`sql.ErrNoRows`这个错误就直接在dao层处理了，向上返一空值, 并且返回是否存在的标识，对于其他错误要`Warp`掉error的， 因为dao层作为应用层的最底层，是最早接触到外部抛出错误的， 保留了本次请求的完整堆栈信息

**dao**

```go
var db *sql.DB

GetUserName(id string)(string, bool, error){
    var result string
    if err := db.QueryRow("select * from user where id = ?", id).Scan(&result);errors.Is(err, sql.ErrNoRows){
       return "", false, nil
    } else {
       return result,err == nil, errors.Wrap(err, "数据查询失败")
    }
}

```

**service**

```go
var ErrNotFind = errors.New("找不到用户")

//  在这个业务中认为找不到用户是一个错误
GetUserName(id string)(string, error){
    if name, ok, err := userDao.GetUserName(id); err != nil {
        return "", err
    } else if !ok{
        return "", ErrNotFind
    } else {
        return name, err
    }
}

```

**controller**

```go
//  api/v1/user/:id
GetUserName(ctx *gin.Context){
    if name,err := userService.GetUserName(id := c.Param("id"));errors.Is(err, userService.ErrNotFind){
        c.String(404, "用户不存在")
    } else if err != nil{
        c.String(500, "服务器正在维护中")
    } else {
        return c.Json(200, map[string]interface{}{
            "Name:"
        })
    }
}
```