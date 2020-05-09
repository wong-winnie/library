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
	conn, err := gorm.Open("mysql", cfg.ConnStr)
	if err != nil {
		log.Fatal("InitGorm Failed", err.Error())
	}
	conn.DB().SetMaxIdleConns(1024)
	conn.DB().SetMaxOpenConns(1024)
	conn.DB().SetConnMaxLifetime(9 * time.Second)
	conn.LogMode(cfg.Debug) //打印SQL
	return &GormMgr{Conn: conn}
}
