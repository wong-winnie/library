package gorm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/wong-winnie/library/dao/common"
	"github.com/wong-winnie/library/dao/config"
	"time"
)

type GormMgr struct {
	Conn *gorm.DB
}

func InitGorm(cfg *config.MysqlCfg) *GormMgr {
	conn, err := gorm.Open("mysql", cfg.ConnStr)
	if err != nil {
		common.SimplePanic("InitGorm失败", err.Error())
	} else {
		conn.DB().SetMaxIdleConns(1024)
		conn.DB().SetMaxOpenConns(1024)
		conn.DB().SetConnMaxLifetime(9 * time.Second)
		conn.LogMode(true) //打印SQL
	}

	return &GormMgr{Conn: conn}
}

func (g *GormMgr) GetSqlDB() *gorm.DB {
	return g.Conn
}
