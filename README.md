# main

```go
package main

import (
	"fmt"
	"github.com/wong-winnie/library/dao/config"
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
}

func InitServiceBlock(cfg *config.Config) *ServiceBlock {
	return &ServiceBlock{
		kafkaProducer: kafka.InitProducer(cfg.KafkaCfg),
		kafkaConsumer: kafka.InitConsumer(cfg.KafkaCfg),
		logConn:       log.InitLog(cfg.LogCfg),
		mysqlConn:     gorm.InitGorm(cfg.MysqlCfg),
		redisConn:     redis.InitRedisSingle(cfg.RedisCfg),
	}
}

func main() {
	blockCfg = &config.Config{
		KafkaCfg: &config.KafkaCfg{Address: []string{"192.168.28.25:9092", "192.168.28.26:9092", "192.168.28.27:9092"}, GroupName: "testKafka6", ClientName: "testKafka-1"},
		LogCfg:   &config.LogCfg{ProgramName: "Test1"},
		MysqlCfg: &config.MysqlCfg{ConnStr: fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "winnie123321", "123.207.79.96:3306", "testDb")},
		RedisCfg: &config.RedisCfg{Addr: "192.168.28.25:6379"},
	}

	blockService = InitServiceBlock(blockCfg)
	fmt.Println("-------InitServiceSuccess-------")

	//TestLog()

	//TestMysql()

	//TestKafkaProducer()

	//TestKafkaConsumer()

	//TestRedis()
}
```

# Log
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

# Kafka
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

#MySQL
```go
package main

import "fmt"

//TODO : StartMySql
func TestMysql() {
	blockService.mysqlConn.Conn.AutoMigrate(TestProgram{})

	data := TestProgram{}
	if ok, err := getOneTestProgram(&data, "id=2"); err != nil && !ok {
		blockService.logConn.ZapErrorLog("getOneTestProgram", "Error", err)
		return
	} else {
		fmt.Println(data)
	}

	data2 := []TestProgram{}
	if ok, err := getAllTestProgram(&data2, "1=1"); err != nil && !ok {
		blockService.logConn.ZapErrorLog("getAllTestProgram", "Error", err)
		return
	} else {
		fmt.Println(data2)
	}

	data3 := []TestProgram{}
	if ok, err := getAllOfPageTestProgram(&data3, "id desc", "1=1", 2, 3); err != nil && !ok {
		blockService.logConn.ZapErrorLog("getAllOfPageTestProgram", "Error", err)
		return
	} else {
		fmt.Println(data3)
	}

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
}

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

//TODO : StartMySql - 查询单条数据
func getOneTestProgram(data *TestProgram, whereSql string) (bool, error) {
	db := blockService.mysqlConn.Conn.Model(data).First(data, whereSql)
	return db.RecordNotFound(), db.Error
}

//TODO : StartMySql - 查询所有数据
func getAllTestProgram(data *[]TestProgram, whereSql string) (bool, error) {
	db := blockService.mysqlConn.Conn.Model(data).Find(data, whereSql)
	return db.RecordNotFound(), db.Error
}

//TODO : StartMySql - 查询所有数据分页
func getAllOfPageTestProgram(data *[]TestProgram, orderSql, whereSql string, page, count int) (bool, error) {
	db := blockService.mysqlConn.Conn.Model(data).Order(orderSql).Offset((page-1)*count).Limit(count).Find(data, whereSql)
	return db.RecordNotFound(), db.Error
}

//TODO : StartMySql - 插入一条或多条数据
//目前 gorm 并不支持批量插入这一功能，但已经被列入 v2.0 的计划里面
func insertTestProgram(data *TestProgram) (int64, error) {
	db := blockService.mysqlConn.Conn.Create(data)
	return db.RowsAffected, db.Error
}

//TODO : StartMySql - 更新一条或多条数据
func updateTestProgram(data *TestProgram, whereSql string) (int64, error) {
	db := blockService.mysqlConn.Conn.Model(data).Where(whereSql).Update(data)
	return db.RowsAffected, db.Error
}

//TODO : StartMySql - 更新一条或多条数据（有零值）
func updateOfZeroTestProgram(data map[string]interface{}, whereSql string) (int64, error) {
	db := blockService.mysqlConn.Conn.Model(TestProgram{}).Where(whereSql).Update(data)
	return db.RowsAffected, db.Error
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