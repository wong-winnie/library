
# 公众号 

![image](https://oscimg.oschina.net/oscnet/up-5fda627460b976ce873642277e9d6e18e72.png)

# main

```go
package main

import (
	"fmt"
	"github.com/wong-winnie/library/dao/config"
	"github.com/wong-winnie/library/dao/elastice"
	"github.com/wong-winnie/library/dao/gorm"
	"github.com/wong-winnie/library/dao/kafka"
	"github.com/wong-winnie/library/dao/log"
	"github.com/wong-winnie/library/dao/redis"
)

var (
	blockCfg     = &config.Config{}
	blockService = &ServiceBlock{}
)

type ServiceBlock struct {
	kafkaProducer *kafka.KafkaProducer
	kafkaConsumer *kafka.KafkaConsumer
	logConn       *log.LogMgr
	mysqlConn     *gorm.GormMgr
	redisConn     *redis.RedisSingleMgr
	elasticConn   *elastice.ElasticMgr
}

func InitServiceBlock(cfg *config.Config) *ServiceBlock {
	return &ServiceBlock{
		kafkaProducer: kafka.InitProducer(cfg.KafkaCfg),
		kafkaConsumer: kafka.InitConsumer(cfg.KafkaCfg),
		logConn:       log.InitLog(cfg.LogCfg),
		mysqlConn:     gorm.InitGorm(cfg.MysqlCfg),
		redisConn:     redis.InitRedisSingle(cfg.RedisCfg),
		elasticConn:   elastice.InitElastic(cfg.ElasticCfg),
	}
}

func main() {
	blockCfg = &config.Config{
		KafkaCfg:   &config.KafkaCfg{Address: []string{"192.168.28.25:9092", "192.168.28.26:9092", "192.168.28.27:9092"}, GroupName: "testKafka6", ClientName: "testKafka-1"},
		LogCfg:     &config.LogCfg{ProgramName: "Test1"},
        //charset=utf8mb4  (utf8不是真正的UTF8,只能包含三个字节的unicode, 4个字节就会报错。改为utf8mb4,就可存入Emoji字符)
		MysqlCfg:   &config.MysqlCfg{ConnStr: fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", "winnie", "winnie", "123.207.79.96:3306", "testDb"), Debug:true},
		RedisCfg:   &config.RedisCfg{Addr: "192.168.28.25:6379"},
		ElasticCfg: &config.ElasticCfg{Url: "http://192.168.28.126:9200"},
	}

	blockService = InitServiceBlock(blockCfg)
	fmt.Println("-------InitServiceSuccess-------")

	//TestLog()

	//TestMysql()

	//TestKafkaProducer()

	//TestKafkaConsumer()

	//TestRedis()

	//TestElastic()
}
```

# patch
```go
package patch

import (
	"fmt"
	"github.com/wong-winnie/library/dao/config"
	"github.com/wong-winnie/library/dao/gorm"
	"github.com/wong-winnie/library/dao/kafka"
	"github.com/wong-winnie/library/dao/log"
	"github.com/wong-winnie/library/dao/redis"
	"sync"
)

var (
	blockCfg     = &config.Config{}
	BlockService = &ServiceBlock{}
)

type ServiceBlock struct {
	KafkaProducer     *kafka.KafkaProducer
	KafkaConsumer     *kafka.KafkaConsumer
	LogConn           *log.LogMgr
	MysqlConn         *gorm.GormMgr
	RedisConn         *redis.RedisSingleMgr
	PoolReqActiveList sync.Pool
	PoolResActiveList sync.Pool
}

func InitServiceBlock(cfg *config.Config) *ServiceBlock {
	return &ServiceBlock{
		KafkaProducer: kafka.InitProducer(cfg.KafkaCfg),
		KafkaConsumer: kafka.InitConsumer(cfg.KafkaCfg),
		LogConn:       log.InitLog(cfg.LogCfg),
		MysqlConn:     gorm.InitGorm(cfg.MysqlCfg),
		RedisConn:     redis.InitRedisSingle(cfg.RedisCfg),
	}
}

func InitPatch() {
	blockCfg = &config.Config{
		KafkaCfg:   &config.KafkaCfg{Address: []string{"192.168.28.25:9092", "192.168.28.26:9092", "192.168.28.27:9092"}, GroupName: "testKafka6", ClientName: "testKafka-1"},
		LogCfg:     &config.LogCfg{ProgramName: "Test1"},
        //charset=utf8mb4  (utf8不是真正的UTF8,只能包含三个字节的unicode, 4个字节就会报错。改为utf8mb4,就可存入Emoji字符)
		MysqlCfg:   &config.MysqlCfg{ConnStr: fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", "winnie", "winnie", "123.207.79.96:3306", "testDb"), Debug: true},
		RedisCfg:   &config.RedisCfg{Addr: "192.168.28.25:6379"},
		ElasticCfg: &config.ElasticCfg{Url: "http://192.168.28.126:9200"},
	}

	BlockService = InitServiceBlock(blockCfg)
	BlockService.LogConn.ZapCustom("LogActiveList")

	fmt.Println("-------InitServiceSuccess-------")
}

```

# Log
//服务运行过程中用vim修改日志数据会导致后续日志无法记录
//vim修改时会创建一个临时文件, 修改完再替换, 原inode被占用, vim用新的inode,替换完变成新的inode
```go
package main

//TODO : StartLog
func TestLog() {
	blockService.logConn.ZapSimpleLog("普通日志", "Key", "Value")
	blockService.logConn.ZapErrorLog("错误日志", "Key", "Value")
	blockService.logConn.ZapPanicLog("异常日志", "Key", "Value")

	//新增日志文件
	blockService.logConn.ZapCustom("TestDemo")
	blockService.logConn.ZapCustomLog("TestDemo", "自定义日志", "Key", "Value")
}
```

# Kafka - [相关文档](https://mp.weixin.qq.com/s?__biz=MzI0OTMzMDgwMg==&mid=2247484202&idx=1&sn=6c3eb1d23f0b8cb7ba2bcfc401ecd6ba)
```go
package main

import "fmt"

func TestKafkaProducer() {
	for i := 0; i < 100; i++ {
		if _, _, err := blockService.kafkaProducer.SendMessage("TestKafkaTopic", []byte(fmt.Sprintf("%d", i))); err != nil {
			blockService.logConn.ZapErrorLog("kafkaProducer", "Error", err)
			continue
		}
	}
}

func TestKafkaConsumer() {
	go func() {
		for err := range blockService.kafkaConsumer.ConsumerConn.Errors() {
			blockService.logConn.ZapErrorLog("kafkaConsumer", "Error", err)
			continue
		}
	}()

	handle := blockService.kafkaConsumer.ConsumerHandler
	handle.RecvFunc = func(recv []byte) bool {
		//TODO : 处理业务
		fmt.Println("Message", string(recv))

		return true
	}

	if err := blockService.kafkaConsumer.TopicConsumer("TestKafkaTopic", &handle); err != nil {
		blockService.logConn.ZapErrorLog("kafkaConsumer", "Error", err)
	}
}

```

# MySQL
```go
package main

import "fmt"

/*
	primary_key 主键
	auto_increment 自增
	not null 非空
	index:索引名 普通索引, 索引名相同->组合索引
	unique_index 唯一索引
	type:text, varchar(255) 数据类型
	default:数值  默认值
*/
type TestProgram struct {
	Id        int    `gorm:"primary_key;auto_increment;not null"`
	Name      string `gorm:"not null;index:sa"`
	Type      int    `gorm:"not null;index:sa"`
	StartTime string `gorm:"not null"`
	AppId     string `gorm:"not null"`
	Secret    string `gorm:"not null"`
}

func TestMysql() {
	//数据库表不存在则创建, 添加缺少的字段，不删除/更改当前数据
	blockService.mysqlConn.Conn.AutoMigrate(TestProgram{})

	//TODO : 查询单条数据
	data := TestProgram{}
	if ok, err := getOneTestProgram(&data, "id=2"); err != nil && !ok {
		blockService.logConn.ZapErrorLog("getOneTestProgram", "Error", err)
		return
	} else {
		fmt.Println(data)
	}

	//TODO : 查询所有数据
	data2 := []TestProgram{}
	if ok, err := getAllTestProgram(&data2, "1=1"); err != nil && !ok {
		blockService.logConn.ZapErrorLog("getAllTestProgram", "Error", err)
		return
	} else {
		fmt.Println(data2)
	}

	//TODO : 查询所有数据分页
	data3 := []TestProgram{}
	if ok, err := getAllOfPageTestProgram(&data3, "id desc", "1=1", 2, 3); err != nil && !ok {
		blockService.logConn.ZapErrorLog("getAllOfPageTestProgram", "Error", err)
		return
	} else {
		fmt.Println(data3)
	}

	//TODO : 插入一条或多条数据
	//目前 gorm 并不支持批量插入这一功能，但已经被列入 v2.0 的计划里面
	data4 := TestProgram{
		Id:        22,
		Name:      "22",
		Type:      22,
		StartTime: "22",
		AppId:     "22",
		Secret:    "22",
	}

	if row, err := insertTestProgram(&data4); err != nil || row == 0 {
		blockService.logConn.ZapErrorLog("insertTestProgram", "Error", err, "Row", row)
		return
	} else {
		fmt.Println(row)
	}

	//TODO : 更新一条或多条数据
	data5 := TestProgram{
		Name: "33",
		Type: 0,
	}
	fmt.Println(data5)
	if row, err := updateTestProgram(&data5, "id=22"); err != nil || row == 0 {
		blockService.logConn.ZapErrorLog("updateTestProgram", "Error", err, "Row", row)
		return
	} else {
		fmt.Println(row)
	}

	//TODO : 更新一条或多条数据（有零值）
	data6 := make(map[string]interface{})
	data6["StartTime"] = "44"
	data6["AppId"] = "55"
	data6["Secret"] = "0"
	if row, err := updateOfZeroTestProgram(data6, "id=22"); err != nil || row == 0 {
		blockService.logConn.ZapErrorLog("updateOfZeroTestProgram", "Error", err, "Row", row)
		return
	} else {
		fmt.Println(row)
	}

	//TODO : 原生sql语句
	data7 := []TestProgram{}
	if err := blockService.mysqlConn.Conn.Raw("select * from test_programs").Scan(&data7).Error; err != nil {
		blockService.logConn.ZapErrorLog("Raw", "Error", err)
		return
	} else {
		fmt.Println(data7)
	}

}

func getOneTestProgram(data *TestProgram, whereSql string) (bool, error) {
	db := blockService.mysqlConn.Conn.Model(data).First(data, whereSql)
	return db.RecordNotFound(), db.Error
}

func getAllTestProgram(data *[]TestProgram, whereSql string) (bool, error) {
	db := blockService.mysqlConn.Conn.Model(data).Find(data, whereSql)
	return db.RecordNotFound(), db.Error
}

func getAllOfPageTestProgram(data *[]TestProgram, orderSql, whereSql string, page, count int) (bool, error) {
	db := blockService.mysqlConn.Conn.Model(data).Order(orderSql).Offset((page-1)*count).Limit(count).Find(data, whereSql)
	return db.RecordNotFound(), db.Error
}

func insertTestProgram(data *TestProgram) (int64, error) {
	db := blockService.mysqlConn.Conn.Create(data)
	return db.RowsAffected, db.Error
}

func updateTestProgram(data *TestProgram, whereSql string) (int64, error) {
	db := blockService.mysqlConn.Conn.Model(data).Where(whereSql).Update(data)
	return db.RowsAffected, db.Error
}

func updateOfZeroTestProgram(data map[string]interface{}, whereSql string) (int64, error) {
	db := blockService.mysqlConn.Conn.Model(TestProgram{}).Where(whereSql).Update(data)
	return db.RowsAffected, db.Error
}
```

# MySQL(v2)
```go
package main

import (
	"errors"
	"fmt"
	"patch"
)

/*
	primary_key 主键
	auto_increment 自增
	not null 非空
	index:索引名 普通索引, 索引名相同->组合索引
	unique_index 唯一索引
	type:text, varchar(255) 数据类型
	default:数值  默认值
*/
type TestProgram struct {
	Id        int    `gorm:"primary_key;auto_increment;not null"`
	Name      string `gorm:"not null;index:sa"`
	Type      int    `gorm:"not null;index:sa"`
	StartTime string `gorm:"not null"`
	AppId     string `gorm:"not null"`
	Secret    string `gorm:"not null"`
}

type ReqTestProgram1 struct {
	Id int64
}

type ResTestProgram1 struct {
	Items TestProgram
}

type ReqTestProgram2 struct {
	pageSize  int
	pageCount int
}

type ResTestProgram2 struct {
	totalCount int64
	Items      []TestProgram
}

type ReqTestProgram3 struct {
	Name      string
	Type      int
	StartTime string
	AppId     string
	Secret    string
}

type ReqTestProgram4 struct {
	Id        int64
	Name      string
	Type      int
	StartTime string
	AppId     string
	Secret    string
}

func main() {
	patch.InitPatch()

	//查询单条数据
	//GetOneData()

	//查询所有数据并分页
	//GetAllData()

	//插入一条或多条数据
	//InsertData()

	//更新一条或多条数据
	//UpdateData()
}

func GetOneData() {
	var (
		req ReqTestProgram1
		res ResTestProgram1
	)

	req = ReqTestProgram1{Id: 3}

	defer func() {
		patch.BlockService.LogConn.ZapErrorLog("GetOneData", "req", req)
		patch.BlockService.LogConn.ZapErrorLog("GetOneData", "res", res)
	}()

	if err := getOneTestProgram(&res.Items, fmt.Sprintf("id=%d", req.Id)); err != nil {
		patch.BlockService.LogConn.ZapErrorLog("GetOneData", "Error", err, "Id", req.Id)
		return
	} else {
		fmt.Println(res)
	}
}

func GetAllData() {
	var (
		req ReqTestProgram2
		res ResTestProgram2
	)

	req = ReqTestProgram2{
		pageSize:  1,
		pageCount: 10,
	}

	defer func() {
		patch.BlockService.LogConn.ZapErrorLog("GetAllData", "req", req)
		patch.BlockService.LogConn.ZapErrorLog("GetAllData", "res", res)
	}()

	if err := GetAllTestProgram(&res.Items, "id desc", "1=1", req.pageSize, req.pageCount, &res.totalCount); err != nil {
		patch.BlockService.LogConn.ZapErrorLog("GetAllData", "Error", err)
		return
	} else {
		fmt.Println(res)
	}
}

func InsertData() {
	var (
		req  ReqTestProgram3
		data TestProgram
	)

	req = ReqTestProgram3{
		Name:      "333",
		Type:      333,
		StartTime: "333",
		AppId:     "333",
		Secret:    "333",
	}

	defer func() {
		patch.BlockService.LogConn.ZapErrorLog("InsertData", "req", req)
	}()

	data = TestProgram{
		Name:      req.Name,
		Type:      req.Type,
		StartTime: req.StartTime,
		AppId:     req.AppId,
		Secret:    req.Secret,
	}

	if err := InsertTestProgram(&data); err != nil {
		patch.BlockService.LogConn.ZapErrorLog("insertTestProgram", "Error", err)
		return
	} else {
		fmt.Println(err)
	}
}

func UpdateData() {
	var (
		req ReqTestProgram4
	)

	req = ReqTestProgram4{
		Id:        333,
		Name:      "",
		Type:      444,
		StartTime: "334443",
		AppId:     "333",
		Secret:    "333",
	}

	defer func() {
		patch.BlockService.LogConn.ZapErrorLog("InsertData", "req", req)
	}()

	data := make(map[string]interface{})
	data["Name"] = req.Name
	data["Type"] = req.Type
	data["AppId"] = req.AppId
	data["Secret"] = req.Secret

	//通过结构体变量更新字段值, gorm库会忽略零值字段。就是字段值等于0, nil, "", false这些值会被忽略掉，不会更新。如果想更新零值，可以使用map类型替代结构体
	if err := UpdateTestProgram(data, fmt.Sprintf("id=%d", req.Id)); err != nil {
		patch.BlockService.LogConn.ZapErrorLog("insertTestProgram", "Error", err)
		return
	} else {
		fmt.Println(err)
	}
}

func getOneTestProgram(data *TestProgram, whereSql string) error {
	//查不到数据时会报错 RecordNotFound
	db := patch.BlockService.MysqlConn.Conn.Model(data).First(data, whereSql)
	return db.Error
}

func GetAllTestProgram(data *[]TestProgram, orderSql, whereSql string, pageSize, pageCount int, totalCount *int64) error {
	db := patch.BlockService.MysqlConn.Conn.Model(data).Where(whereSql).Count(totalCount).Order(orderSql).Offset((pageSize - 1) * pageCount).Limit(pageCount).Find(data)
	return db.Error
}

func InsertTestProgram(data *TestProgram) error {
	db := patch.BlockService.MysqlConn.Conn.Create(data)
	if db.RowsAffected == 0 {
		return errors.New("RowsAffected Is Zero")
	}
	return db.Error
}

func UpdateTestProgram(data map[string]interface{}, whereSql string) error {
	db := patch.BlockService.MysqlConn.Conn.Model(TestProgram{}).Where(whereSql).Update(data)
	if db.RowsAffected == 0 {
		return errors.New("RowsAffected Is Zero")
	}
	return db.Error
}

```

# Redis
```go
package main

import "fmt"

func TestRedis(){
	if err :=blockService.redisConn.Conn.Set("test", 1, 0).Err(); err != nil{
		blockService.logConn.ZapErrorLog("redisConn", "Error", err)
		return
	}

	v :=blockService.redisConn.Conn.Get("test").String()
	fmt.Println(v)
}
```

# ElasticSearch
```go
package main

import (
	"context"
	"fmt"
)

type Employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

func TestElastic() {
	e := Employee{"winnie", "wong", 18, "good good study, day day up", []string{"IT"}}
	put, err := blockService.elasticConn.Conn.Index().
		Index("employee_index").
		Type("employee_type").
		Id("2").
		BodyJson(e).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(put)

	result, err := blockService.elasticConn.Conn.Search().Index("employee_index").Type("employee_type").Do(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}

	for k, v := range result.Hits.Hits {
		fmt.Println(k, string(v.Source))
	}
}

```