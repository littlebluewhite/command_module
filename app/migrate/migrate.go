package migrate

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"new_command/config"
)

func NewMigrateDB() (*sql.DB, error) {
	DBConfig := config.NewConfig[config.DBConfig](".", "env", "db")
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", DBConfig.Host, DBConfig.Port),
		DBName:               DBConfig.DB,
		User:                 DBConfig.User,
		Passwd:               DBConfig.Password,
		AllowNativePasswords: false,
		ParseTime:            true,
	}

	return sql.Open("mysql", cfg.FormatDSN())
}
