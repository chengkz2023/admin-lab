package request

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
)

type InitDB struct {
	AdminPassword string `json:"adminPassword" binding:"required"`
	Host          string `json:"host"`
	Port          string `json:"port"`
	UserName      string `json:"userName"`
	Password      string `json:"password"`
	DBName        string `json:"dbName" binding:"required"`
}

func (i *InitDB) MysqlEmptyDsn() string {
	if i.Host == "" {
		i.Host = "127.0.0.1"
	}
	if i.Port == "" {
		i.Port = "3306"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/", i.UserName, i.Password, i.Host, i.Port)
}

func (i *InitDB) ToMysqlConfig() config.Mysql {
	return config.Mysql{
		GeneralDB: config.GeneralDB{
			Path:         i.Host,
			Port:         i.Port,
			Dbname:       i.DBName,
			Username:     i.UserName,
			Password:     i.Password,
			MaxIdleConns: 10,
			MaxOpenConns: 100,
			LogMode:      "error",
			Config:       "charset=utf8mb4&parseTime=True&loc=Local",
		},
	}
}
