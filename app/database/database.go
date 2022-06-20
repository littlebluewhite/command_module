package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"new_command/config"
	"path"
	"path/filepath"
	"runtime"
)

var (
	DBConfig *config.DBConfig
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	DBConfig = config.LoadConfig[*(config.DBConfig)](path.Join(rootPath, "env"), "db")

}

func NewDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&parseTime=true&loc=Local",
		DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.DB)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func CloseDB(*gorm.DB) error {
	//Todo Check the error
	return nil
}
