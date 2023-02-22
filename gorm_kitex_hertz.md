# Gorm_Kitex_Hertz

## Gorm
Gorm支持的数据库有Mysql,SqlServer,PostagreSql,Sqlite

使用时需要导入mysql引擎`"gorm.io/driver/mysql"`

[DSN具体内容](https://github.com/go-sql-driver/mysql#dsn-data-source-name)
### Connect 
```GO
type Product struct {
	gorm.Model
	Code  string
	Price uint
}
//其中gorm.Model提供主键、创建时间、销毁时间等信息
dsn := "root:@tcp(127.0.0.1:3306)/gorm?&charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err != nil {
//log.Fatal(err)
panic(err)
}
```
### Create
```Go
//插入单条数据
db.Create(&Product{Code: "D42", Price: 100})

//插入多条数据
projects := []*Project{{Code: "D40"}, {Code: "D41"}, {Code: "D42"}}
res = db.Create(projects)
fmt.Println(res.Error)
for _, p := range projects {
fmt.Println(p.ID)
}
//处理冲突,如何使用Upsert

```
如果需要默认值可以使用，gorm的反射在结构体变量中加入`gorm:"deafault:19" `
如果想改变字段值，可以使用`gorm:"column:code"` 
![](https://cdn.jsdelivr.net/gh/richardzhangy26/Pic@main/src202301231438920.png)
可以看`gorm.Model`提供主键Id、创建时间、销毁时间等信息
Gorm 的约定（默认）
- Gorm 使用名为ID 的字段作为主键
- 使用结构体的蛇形负数作为表名
- 字段名的蛇形作为列名
- 使用 CreatedAt、 UpdatedAt 字段作为创建、 更新时间

### Read
```Go
var product Product 
//单条查询
db.First(&product, 1)                 // 根据整形主键查找
db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
result := db.First(&user)
result.RowsAffected // 返回找到的记录数
result.Error        // returns error

//多条查询
db.Where("name <> ?", "jinzhu").Find(&users)
// SELECT * FROM users WHERE name <> 'jinzhu';

// IN
db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
// SELECT * FROM users WHERE name IN ('jinzhu','jinzhu 2');

// LIKE
db.Where("name LIKE ?", "%jin%").Find(&users)
// SELECT * FROM users WHERE name LIKE '%jin%';

// AND
db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;
// Struct
db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;

// Map
db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;

// 主键切片条件
db.Where([]int64{20, 21, 22}).Find(&users)
// SELECT * FROM users WHERE id IN (20, 21, 22);
```
值得注意的是

使用 First 时，需要注意查询不到数据会返回 ErrRecordNotFound。

使用 Find 查询多条数据，查询不到数据不会返回错误。

使用结构体作为查询条件

当使用结构作为条件查询时，GORM 只会查询非零值字段。这意味着如果您的字段值为 O、"、false 或其他零值， 该字段不会被用于构建查询条件，使用 Map 来构建查询条件。
### Update
```Go
db.Model(&product).Update("Price", 200)
// Update - 更新多个字段
db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

// 根据 `struct` 更新属性，只会更新非零值的字段
db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

// 根据 `map` 更新属性
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
// UPDATE users SET name='hello', age=18, actived=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
```
注意 当通过 struct 更新时，GORM 只会更新非零字段。 如果您想确保指定字段被更新，你应该使用 Select 更新选定字段，或使用 map 来完成更新操作

### Delete
```Go
//物理删除
db.Delete(&product,1)
//软删除
// user 的 ID 是 `111`
db.Delete(&user)
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;

// 批量删除
db.Where("age = ?", 20).Delete(&User{})
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

// 在查询时会忽略被软删除的记录
db.Where("age = 20").Find(&user)
// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
//查找被软删除的记录
//您可以使用 Unscoped 找到被软删除的记录

db.Unscoped().Where("age = 20").Find(&users)
// SELECT * FROM users WHERE age = 20;
//永久删除
//您也可以使用 Unscoped 永久删除匹配的记录

db.Unscoped().Delete(&order)
// DELETE FROM orders WHERE id=10;

```
如果您的模型包含了一个 gorm.deletedat 字段（gorm.Model 已经包含了该字段)，它将自动获得软删除的能力！
拥有软删除能力的模型调用 Delete 时，记录不会被从数据库中真正删除。但 GORM 会将 DeletedAt 置为当前时间， 并且你不能再通过正常的查询方法找到该记录。

### 事务

为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，您可以在初始化时禁用它，这将获得大约 30%+ 性能提升。
`// 全局禁用
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
SkipDefaultTransaction: true,
})`
```go
db.Transaction(func(tx *gorm.DB) error {
  // 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
  if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
    // 返回任何错误都会回滚事务
    return err
  }

  if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
    return err
  }

  // 返回 nil 提交事务
  return nil
})
//手动执行
// 开始事务
tx := db.Begin()

// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
tx.Create(...)

// ...

// 遇到错误时回滚事务
tx.Rollback()

// 否则，提交事务
tx.Commit()
```
为了避免忘记commit,rollback还是需要使用transaction来自动执行事务

### Hook

Hook 是在创建、查询、更新、删除等操作之前、之后调用的函数。

如果您已经为模型定义了指定的方法，它会在创建、更新、查询、删除时自动被调用。如果任何回调返回错误，GORM 将停止后续的操作并回滚事务。

钩子方法的函数签名应该是 func(*gorm.DB) error

```go
// 开始事务
BeforeSave
BeforeCreate
// 关联前的 save
// 插入记录至 db
// 关联后的 save
AfterCreate
AfterSave
// 提交或回滚事务

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
u.UUID = uuid.New()

if !u.IsValid() {
err = errors.New("can't save invalid data")
}
return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
if u.ID == 1 {
tx.Model(u).Update("role", "admin")
}
return
} 
```

注意 在 GORM 中保存、删除操作会默认运行在事务上， 因此在事务完成之前该事务中所作的更改是不可见的，如果您的钩子返回了任何错误，则修改将被回滚。


## Kitex & Hertz初体验

Kitex是字节跳动内部的Golang微服务RPC框架，具有**高性能、强可扩展**的特点

Kitex 框架及命令行工具，默认支持 thrift 和 proto3 两种 IDL，对应的 Kitex 支持 thrift 和 protobuf 两种序列化协议。 传输上 Kitex 使用扩展的 thrift 作为底层的传输协议（注：thrift 既是 IDL 格式，同时也是序列化协议和传输协议）。

### IDL
IDL 全称是 Interface Definition Language，接口定义语言

是用来描述软件组件接口的一种计算机语言。IDL通过一种独立于编程语言的方式来描述接口，使得在不同平台上运行的对象和用不同语言编写的程序可以相互通信交流；


｜消息类型	｜编码协议	｜传输协议｜

｜PingPong｜Thrift / Protobuf｜	TTHeader / HTTP2(gRPC)

｜Oneway｜Thrift	      ｜TTHeader｜

｜Streaming｜	Protobuf｜	HTTP2(gRPC)｜

### 例子
先进入hello目录
`cd kitex-examples/hello`
运行server
`go run .`
运行client
`go run ./client`
![](https://cdn.jsdelivr.net/gh/richardzhangy26/Pic@main/src202301271651490.png)

### Thrift
目前 Thrift 支持 PingPong 和 Oneway。Kitex 计划支持 Thrift Streaming。

其IDL定义
```Thrift
namespace go echo

struct Request {
    1: string Msg
}

struct Response {
    1: string Msg
}

service EchoService {
    Response Echo(1: Request req
    ); // pingpong method
    oneway void VisitOneway(1: Request req); // oneway method
}

```
### Protobuf
itex 支持两种承载 Protobuf 负载的协议：

Kitex Protobuf
- 只支持 PingPong，若 IDL 定义了 stream 方法，将默认使用 gRPC 协议
gRPC 协议
- 可以与 gRPC 互通，与 gRPC service 定义相同，支持 Unary(PingPong)、 Streaming 调用

```Protobuf
syntax = "proto3";

option go_package = "echo";

package echo;

message Request {
  string msg = 1;
}

message Response {
  string msg = 1;
}

service EchoService {
  rpc ClientSideStreaming(stream Request) returns (Response) {} // 客户端侧 streaming
  rpc ServerSideStreaming(Request) returns (stream Response) {} // 服务端侧 streaming
  rpc BidiSideStreaming(stream Request) returns (stream Response) {} // 双向流
}

```

### kitex生成工具

当你修改完`.thrift`文件后可以使用

```Go
kitex -service a.b.c hello.thrift

# 若当前目录不在 $GOPATH/src 下，需要加上 -module 参数，一般为 go.mod 下的名字
kitex -module "your_module_name" -service a.b.c hello.thrift

```
kitex会更新 `./handler.go`，在里面增加一个 Add 方法的基本实现
更新 `./kitex_gen`，里面有框架运行所必须的代码文件

### kitex服务注册与发现
Kitex 已经通过社区开发者的支持，完成了 ETCD、ZooKeeper、Eureka、Consul、Nacos、Polaris 多种服务发现模式，当然也支持 DNS 解析以及 Static IP 直连访问模式，建立起了强大且完备的社区生态，供用户按需灵活选用

etcd（读作 et-see-dee）是一种开源的分布式统一键值存储，用于分布式系统或计算机集群的共享配置、服务发现和的调度协调。etcd 有助于促进更加安全的自动更新，协调向主机调度的工作，并帮助设置容器的覆盖网络。
```Go
package main

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api/hello"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	client := hello.MustNewClient("Hello", client.WithResolver(r))
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		resp, err := client.Echo(ctx, &api.Request{Message: "Hello"})
		cancel()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
		time.Sleep(time.Second)
	}
}

//https://github.com/kitex-contrib/registry-etcd/blob/main/example/client/main.go
```
# GO语言笔记服务

