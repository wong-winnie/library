package gorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"library/dao/common"
	"library/dao/config"
	"time"
)

type GormMgr struct {
	Conn *gorm.DB
}

func InitGorm(cfg *config.MysqlCfg) *GormMgr {
	conn, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Address, cfg.DBName))
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

func InitGorm2(cfg *config.MysqlCfg) *GormMgr {
	conn, err := gorm.Open("mysql", cfg.DSN)
	if err != nil {
		common.SimplePanic("InitGorm失败", err.Error())
	} else {
		conn.DB().SetMaxIdleConns(1024)
		conn.DB().SetMaxOpenConns(1024)
		conn.DB().SetConnMaxLifetime(9 * time.Second)

		//conn.LogMode(true)  //打印SQL
	}

	return &GormMgr{Conn: conn}
}

func (g *GormMgr) GetSqlDB() *gorm.DB {
	return g.Conn
}
