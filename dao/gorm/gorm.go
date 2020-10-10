package gorm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/wong-winnie/library/dao/config"
	"log"
	"time"
)

type GormMgr struct {
	Conn *gorm.DB
}

func InitGorm(cfg *config.MysqlCfg) *GormMgr {
	//用户名或密码有特殊字符（例如@、#）时就会报错，解决方法就是将特殊字符进行转义，方法如下：
	//user := url.QueryEscape(user)
	//password := url.QueryEscape(password)

	conn, err := gorm.Open(cfg.DriveType, cfg.ConnStr)
	if err != nil {
		log.Fatal("InitGorm Failed", err.Error())
	}
	//配置数偏少于MySQL配置的最大连接数(show variables like '%max_connections%')
	conn.DB().SetMaxIdleConns(1024)
	conn.DB().SetMaxOpenConns(1024)
	conn.DB().SetConnMaxLifetime(9 * time.Second)
	conn.LogMode(cfg.Debug) //打印SQL
	return &GormMgr{Conn: conn}
}
