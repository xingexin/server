package db

import (
	"fmt"
	serverConfig "server/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {

	cfg, err := serverConfig.LoadConfig()
	if err != nil {
		return nil, err
	}
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
