package db

import (
	"fmt"
	"server/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接，支持 MySQL 和 PostgreSQL
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var diaLector gorm.Dialector
	switch cfg.DataBase.Driver {
	case "mysql":
		diaLector = mysql.Open(cfg.DataBase.DSN)
	case "pgsql":
		diaLector = postgres.Open(cfg.DataBase.DSN)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.DataBase.Driver)
	}
	DB, err := gorm.Open(diaLector)
	if err != nil {
		panic("fail to connect database")
	}
	return DB, nil
}
